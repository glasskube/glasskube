package web

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/dependency/adapter/goclient"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/glasskube/glasskube/pkg/update"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
	"github.com/glasskube/glasskube/internal/web/components/pkg_update_alert"
	"github.com/glasskube/glasskube/internal/web/handler"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/manifest"
	"github.com/glasskube/glasskube/pkg/open"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/uninstall"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

//go:embed root
//go:embed templates
var embededFs embed.FS

type ServerOptions struct {
	Host       string
	Port       int32
	Kubeconfig string
}

func NewServer(options ServerOptions) *server {
	server := server{
		ServerOptions: options,
		configLoader:  &defaultConfigLoader{options.Kubeconfig},
		forwarders:    make(map[string]*open.OpenResult),
	}
	return &server
}

type server struct {
	ServerOptions
	configLoader
	listener      net.Listener
	restConfig    *rest.Config
	rawConfig     *api.Config
	pkgClient     client.PackageV1Alpha1Client
	wsHub         *WsHub
	informerStore cache.Store
	informerCtrl  cache.Controller
	forwarders    map[string]*open.OpenResult
	dependencyMgr *dependency.DependendcyManager
}

func (s *server) RestConfig() *rest.Config {
	return s.restConfig
}

func (s *server) RawConfig() *api.Config {
	return s.rawConfig
}

func (s *server) Client() client.PackageV1Alpha1Client {
	return s.pkgClient
}

func (s *server) broadcastPkg(pkg *v1alpha1.Package, status *client.PackageStatus, installedManifest *v1alpha1.PackageManifest, latestVersion string) {
	go func() {
		var bf bytes.Buffer
		err := pkg_overview_btn.Render(&bf, pkgOverviewBtnTmpl, pkg, status, installedManifest, latestVersion)
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_overview_btn.TemplateId, pkg.Name))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()
	go func() {
		var bf bytes.Buffer
		err := pkg_detail_btns.Render(&bf, pkgDetailBtnsTmpl, pkg, status, installedManifest, latestVersion)
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_detail_btns.TemplateId, pkg.Name))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()
	go func() {
		var bf bytes.Buffer
		err := pkg_update_alert.Render(&bf, pkgUpdateAlertTmpl, s.isUpdateAvailable(context.TODO()))
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_update_alert.TemplateId, pkg.Name))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()
}

func (s *server) Start(ctx context.Context) error {
	if s.listener != nil {
		return errors.New("server is already listening")
	}

	_ = s.initKubeConfig()
	if err := s.checkBootstrapped(ctx); err == nil {
		s.startInformer(ctx)
	}

	root, err := fs.Sub(embededFs, "root")
	if err != nil {
		return err
	}

	s.wsHub = NewHub()
	fileServer := http.FileServer(http.FS(root))

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(fileServer)
	router.Handle("/favicon.ico", fileServer)
	router.HandleFunc("/ws", s.wsHub.handler)
	router.HandleFunc("/support", s.supportPage)
	router.HandleFunc("/kubeconfig", s.kubeconfigPage)
	router.Handle("/bootstrap", s.requireKubeconfig(s.bootstrapPage))
	router.Handle("/kubeconfig/persist", s.requireKubeconfig(s.persistKubeconfig))
	router.Handle("/packages", s.requireReady(s.packages))
	router.Handle("/packages/install", s.requireReady(s.install))
	router.Handle("/packages/install/modal", s.requireReady(s.installModal))
	router.Handle("/packages/update", s.requireReady(s.update))
	router.Handle("/packages/update/modal", s.requireReady(s.updateModal))
	router.Handle("/packages/uninstall", s.requireReady(s.uninstall))
	router.Handle("/packages/open", s.requireReady(s.open))
	router.Handle("/packages/{pkgName}", s.requireReady(s.packageDetail))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/packages", http.StatusFound) })
	http.Handle("/", s.enrichContext(router))

	bindAddr := fmt.Sprintf("%v:%d", s.Host, s.Port)

	s.listener, err = net.Listen("tcp", bindAddr)
	if err != nil {
		// Checks if Port Conflict Error exists
		if isPortConflictError(err) {
			userInput := cliutils.YesNoPrompt(
				"Port is already in use.\nShould glasskube use a different port? (Y/n): ", true)
			if userInput {
				s.listener, err = net.Listen("tcp", ":0")
				if err != nil {
					panic(err)
				}
				bindAddr = fmt.Sprintf("%v:%d", s.Host, s.listener.Addr().(*net.TCPAddr).Port)
			} else {
				fmt.Println("Exiting. User chose not to use a different port.")
				os.Exit(1)
			}
		} else {
			// If no Port Conflict error is found, return other errors
			return err
		}
	}

	fmt.Printf("glasskube UI is available at http://%v\n", bindAddr)
	_ = cliutils.OpenInBrowser("http://" + bindAddr)

	go s.wsHub.Run()
	server := &http.Server{}
	err = server.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *server) install(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	selectedVersion := r.FormValue("selectedVersion")
	enableAutoUpdateVal := r.FormValue("enableAutoUpdate")
	if selectedVersion == "" {
		var packageIndex repo.PackageIndex
		if err := repo.FetchPackageIndex("", pkgName, &packageIndex); err != nil {
			fmt.Fprintf(os.Stderr, "❗ Error: Could not fetch package metadata: %v\n", err)
			return
		}
		selectedVersion = packageIndex.LatestVersion
	}
	err := install.NewInstaller(s.pkgClient).
		WithStatusWriter(statuswriter.Stderr()).
		Install(r.Context(), pkgName, selectedVersion, strings.ToLower(enableAutoUpdateVal) == "on")
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred installing %v: \n%v\n", pkgName, err)
	}
}

func (s *server) installModal(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	selectedVersion := r.FormValue("selectedVersion")
	var idx repo.PackageIndex
	if err := repo.FetchPackageIndex("", pkgName, &idx); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred fetching versions of %v: %v\n", pkgName, err)
		return
	}
	if selectedVersion == "" {
		selectedVersion = idx.LatestVersion
	}
	var mf v1alpha1.PackageManifest
	if err := repo.FetchPackageManifest("", pkgName, selectedVersion, &mf); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred fetching manifest of %v in version %v: %v\n", pkgName, selectedVersion, err)
		return
	}

	res, err := s.dependencyMgr.Validate(r.Context(), &mf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred validating dependencies of %v in version %v: %v\n", pkgName, selectedVersion, err)
		return
	}

	err = pkgInstallModalTmpl.Execute(w, &map[string]any{
		"PackageName":      pkgName,
		"PackageIndex":     &idx,
		"SelectedVersion":  selectedVersion,
		"ShowConflicts":    res.Status == dependency.ValidationResultStatusConflict,
		"ValidationResult": res,
	})
	checkTmplError(err, "pkgInstallModalTmpl")
}

func (s *server) updateModal(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	pkgs := make([]string, 0, 1)
	if pkgName != "" {
		pkgs = append(pkgs, pkgName)
	}

	updates := make([]*map[string]any, 0)
	updater := update.NewUpdater(s.pkgClient).WithStatusWriter(statuswriter.Stderr())
	ut, err := updater.Prepare(r.Context(), pkgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred preparing update of %v: \n%v\n", pkgName, err)
		return
	}
	for _, u := range ut.Items {
		if u.UpdateRequired() {
			updates = append(updates, &map[string]any{
				"Name":           u.Package.Name,
				"CurrentVersion": u.Package.Spec.PackageInfo.Version,
				"LatestVersion":  u.Version,
			})
		}
	}

	err = pkgUpdateModalTmpl.Execute(w, &map[string]any{
		"Updates":     updates,
		"PackageName": pkgName,
	})
	checkTmplError(err, "pkgUpdateModalTmpl")
}

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := r.FormValue("packageName")
	pkgs := make([]string, 0, 1)
	if pkgName != "" {
		pkgs = append(pkgs, pkgName)
	}

	updater := update.NewUpdater(s.pkgClient).WithStatusWriter(statuswriter.Stderr())
	ut, err := updater.Prepare(ctx, pkgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred preparing update of %v: \n%v\n", pkgName, err)
		return
	}
	// in the future we might want to check here whether the prepared new version is the same as the "toVersion"
	// which the user agreed to update to in the dialog
	err = updater.Apply(ctx, ut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred updating %v: \n%v\n", pkgName, err)
		return
	}
}

func (s *server) uninstall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := r.FormValue("packageName")
	var pkg v1alpha1.Package
	if err := s.pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred fetching %v during uninstall: \n%v\n", pkgName, err)
		return
	}
	if err := uninstall.NewUninstaller(s.pkgClient).
		WithStatusWriter(statuswriter.Stderr()).
		Uninstall(ctx, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
	}
}

func (s *server) open(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	if result, ok := s.forwarders[pkgName]; ok {
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		return
	}

	result, err := open.NewOpener().Open(r.Context(), pkgName, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open %v: %v\n", pkgName, err)
	} else {
		s.forwarders[pkgName] = result
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *server) packages(w http.ResponseWriter, r *http.Request) {
	packages, err := list.GetPackagesWithStatus(s.pkgClient, r.Context(), list.ListOptions{IncludePackageInfos: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load packages: %v\n", err)
		return
	}

	err = pkgsPageTmpl.Execute(w, &map[string]any{
		"CurrentContext":   s.rawConfig.CurrentContext,
		"Packages":         packages,
		"UpdatesAvailable": s.isUpdateAvailable(r.Context()),
	})
	checkTmplError(err, "packages")
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkgName"]
	pkg, status, manifest, latestVersion, err := describe.DescribePackage(r.Context(), pkgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred fetching package details of %v: \n%v\n", pkgName, err)
		return
	}
	if latestVersion == "" {
		// TODO have a look at handling latestVersion as return value from DescribePackage again – seems weird
		latestVersion, err = repo.GetLatestVersion("", pkgName)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred fetching latest version of %v: \n%v\n", pkgName, err)
	}
	err = pkgPageTmpl.Execute(w, &map[string]any{
		"CurrentContext": s.rawConfig.CurrentContext,
		"Package":        pkg,
		"Status":         status,
		"Manifest":       manifest,
		"LatestVersion":  latestVersion,
	})
	checkTmplError(err, fmt.Sprintf("package-detail (%s)", pkgName))
}

func (s *server) supportPage(w http.ResponseWriter, r *http.Request) {
	if err := s.checkBootstrapped(r.Context()); err != nil {
		if err.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
			return
		}
		err := supportPageTmpl.Execute(w, &map[string]any{
			"CurrentContext":            "",
			"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
			"Err":                       err,
		})
		checkTmplError(err, "support")
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *server) bootstrapPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == "POST" {
		client := bootstrap.NewBootstrapClient(s.restConfig)
		if err := client.Bootstrap(ctx, bootstrap.DefaultOptions()); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			err := bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-failure", nil)
			checkTmplError(err, "bootstrap-failure")
		} else {
			err := bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-success", nil)
			checkTmplError(err, "bootstrap-success")
		}
	} else {
		isBootstrapped, err := bootstrap.IsBootstrapped(ctx, s.restConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
		} else if isBootstrapped {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		tplErr := bootstrapPageTmpl.Execute(w, &map[string]any{
			"CurrentContext": s.rawConfig.CurrentContext,
			"Err":            err,
		})
		checkTmplError(tplErr, "bootstrap")
	}
}

func (s *server) kubeconfigPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("kubeconfig")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.loadBytesConfig(data)
		if err := s.checkKubeconfig(); err != nil {
			fmt.Fprintf(os.Stderr, "The selected kubeconfig is invalid: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "The selected kubeconfig is valid!")
		}
	}

	configErr := s.checkKubeconfig()
	var currentContext string
	if s.rawConfig != nil {
		currentContext = s.rawConfig.CurrentContext
	}
	tplErr := kubeconfigPageTmpl.Execute(w, map[string]any{
		"CurrentContext":            currentContext,
		"ConfigErr":                 configErr,
		"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
		"DefaultKubeconfigExists":   defaultKubeconfigExists(),
	})
	checkTmplError(tplErr, "kubeconfig")
}

func (s *server) persistKubeconfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !defaultKubeconfigExists() {
			if err := clientcmd.WriteToFile(*s.rawConfig, clientcmd.RecommendedHomeFile); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Redirect(w, r, "/", http.StatusFound)
			}
		} else {
			fmt.Fprintln(os.Stderr, "default kubeconfig already exists! nothing was saved")
			http.Error(w, "", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "only POST is supported", http.StatusMethodNotAllowed)
	}
}

func (server *server) loadBytesConfig(data []byte) {
	server.configLoader = &bytesConfigLoader{data}
}

func (server *server) checkKubeconfig() ServerConfigError {
	if server.pkgClient == nil {
		return server.initKubeConfig()
	} else {
		return nil
	}
}

func (server *server) checkBootstrapped(ctx context.Context) ServerConfigError {
	if err := server.checkKubeconfig(); err != nil {
		return err
	}

	isBootstrapped, err := bootstrap.IsBootstrapped(ctx, server.restConfig)
	if !isBootstrapped || err != nil {
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
		}
		return newBootstrapErr(err)
	}

	return nil
}

func (server *server) initKubeConfig() ServerConfigError {
	restConfig, rawConfig, err := server.LoadConfig()
	if err != nil {
		return newKubeconfigErr(err)
	}
	client, err := client.New(restConfig)
	if err != nil {
		return newKubeconfigErr(err)
	}
	server.restConfig = restConfig
	server.rawConfig = rawConfig
	server.pkgClient = client
	server.dependencyMgr = dependency.NewDependencyManager(goclient.NewGoClientAdapter(server.pkgClient), owners.NewOwnerManager(scheme.Scheme))
	return nil
}

func (server *server) startInformer(ctx context.Context) {
	if server.informerStore == nil && server.informerCtrl == nil {
		server.informerStore, server.informerCtrl = server.initInformer(ctx)
		go server.informerCtrl.Run(ctx.Done())
		server.pkgClient = server.pkgClient.WithPackageStore(server.informerStore)
	}
}

func (s *server) enrichContext(h http.Handler) http.Handler {
	return &handler.ContextEnrichingHandler{Source: s, Handler: h}
}

func (s *server) requireReady(h http.HandlerFunc) http.Handler {
	return &handler.PreconditionHandler{
		Precondition: func(r *http.Request) error {
			err := s.checkBootstrapped(r.Context())
			if err != nil {
				return err
			}
			s.startInformer(context.WithoutCancel(r.Context()))
			return nil
		},
		Handler:       h,
		FailedHandler: handleConfigError,
	}
}

func (s *server) requireKubeconfig(h http.HandlerFunc) http.Handler {
	return &handler.PreconditionHandler{
		Precondition:  func(r *http.Request) error { return s.checkKubeconfig() },
		Handler:       h,
		FailedHandler: handleConfigError,
	}
}

func handleConfigError(w http.ResponseWriter, r *http.Request, err error) {
	if sce, ok := err.(ServerConfigError); ok {
		if sce.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
		} else {
			http.Redirect(w, r, "/support", http.StatusFound)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func defaultKubeconfigExists() bool {
	if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil {
		return true
	} else {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "could not check kubeconfig file: %v", err)
		}
		return false
	}
}

func (s *server) initInformer(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.pkgClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var pkgList v1alpha1.PackageList
				err := pkgClient.Packages().GetAll(ctx, &pkgList)
				return &pkgList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.Packages().Watch(ctx)
			},
		},
		&v1alpha1.Package{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcastPkg(pkg, client.GetStatusOrPending(&pkg.Status), nil, "")
				}
			},
			UpdateFunc: func(oldObj, newObj any) {
				if pkg, ok := newObj.(*v1alpha1.Package); ok {
					ctx := client.SetupContextWithClient(ctx, s.restConfig, s.rawConfig, s.pkgClient)
					if mf, err := manifest.GetInstalledManifestForPackage(ctx, *pkg); err != nil && !errors.Is(err, manifest.ErrPackageNoManifest) {
						fmt.Fprintf(os.Stderr, "Error fetching manifest for package %v: %v\n", pkg.Name, err)
					} else {
						if latestVersion, err := repo.GetLatestVersion("", pkg.Name); err != nil {
							fmt.Fprintf(os.Stderr, "An error occurred fetching latest version of %v: \n%v\n", pkg.Name, err)
						} else {
							s.broadcastPkg(pkg, client.GetStatusOrPending(&pkg.Status), mf, latestVersion)
						}
					}
				}
			},
			DeleteFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcastPkg(pkg, client.GetStatus(&pkg.Status), nil, "")
				}
			},
		},
	)
}

func (s *server) isUpdateAvailable(ctx context.Context) bool {
	if tx, err := update.NewUpdater(s.pkgClient).Prepare(ctx, []string{}); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v", err)
		return false
	} else {
		return !tx.IsEmpty()
	}
}

func isPortConflictError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if osErr, ok := opErr.Err.(*os.SyscallError); ok {
			return osErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}
