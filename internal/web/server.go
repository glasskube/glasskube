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
	"syscall"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
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
	listener   net.Listener
	restConfig *rest.Config
	rawConfig  *api.Config
	pkgClient  *client.PackageV1Alpha1Client
	wsHub      *WsHub
	forwarders map[string]*open.OpenResult
}

func (s *server) RestConfig() *rest.Config {
	return s.restConfig
}

func (s *server) RawConfig() *api.Config {
	return s.rawConfig
}

func (s *server) Client() *client.PackageV1Alpha1Client {
	return s.pkgClient
}

func (s *server) broadcastPkgStatusUpdate(
	pkgName string,
	status *client.PackageStatus,
	manifest *v1alpha1.PackageManifest,
) {
	go func() {
		var bf bytes.Buffer
		err := pkg_overview_btn.Render(&bf, pkgOverviewBtnTmpl, pkgName, status, manifest)
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_overview_btn.TemplateId, pkgName))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()
	go func() {
		var bf bytes.Buffer
		err := pkg_detail_btns.Render(&bf, pkgDetailBtnsTmpl, pkgName, status, manifest)
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_overview_btn.TemplateId, pkgName))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()
}

func (s *server) Start() error {
	if s.listener != nil {
		return errors.New("server is already listening")
	}

	_ = s.initKubeConfig()

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
	router.Handle("/packages", s.requireBootstrapped(s.packages))
	router.Handle("/packages/install", s.requireBootstrapped(s.install))
	router.Handle("/packages/install/modal", s.requireBootstrapped(s.installModal))
	router.Handle("/packages/install/modal/versions", s.requireBootstrapped(s.installModalVersions))
	router.Handle("/packages/uninstall", s.requireBootstrapped(s.uninstall))
	router.Handle("/packages/open", s.requireBootstrapped(s.open))
	router.Handle("/packages/{pkgName}", s.requireBootstrapped(s.packageDetail))
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
	ctxAsync := context.WithoutCancel(r.Context())
	go func() {
		status, err := install.NewInstaller(s.pkgClient).
			WithStatusWriter(statuswriter.Stderr()).
			InstallBlocking(ctxAsync, pkgName, selectedVersion)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred installing %v: \n%v\n", pkgName, err)
			return
		}
		manifest, err := manifest.GetInstalledManifest(ctxAsync, pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not fetch manifest of %v: %v\n", pkgName, err)
			return
		}
		s.broadcastPkgStatusUpdate(pkgName, status, manifest)
	}()
	s.broadcastPkgStatusUpdate(pkgName, client.NewPendingStatus(), nil)
}

func (s *server) installModal(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")

	_, _, manifest, err := describe.DescribePackage(r.Context(), pkgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred fetching package details of %v: \n%v\n", pkgName, err)
		return
	}

	err = pkgInstallModalTmpl.Execute(w, &map[string]any{
		"Manifest": manifest,
	})
	checkTmplError(err, "pkgInstallModalTmpl")
}

func (s *server) installModalVersions(w http.ResponseWriter, r *http.Request) {
	enableAutoUpdateVal := r.FormValue("enableAutoUpdate")
	showVersions := enableAutoUpdateVal == ""
	var idx repo.PackageIndex
	if showVersions {
		pkgName := r.FormValue("packageName")
		if err := repo.FetchPackageIndex("", pkgName, &idx); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred fetching versions of %v: %v\n", pkgName, err)
		}
	}

	err := pkgInstallModalVersionsTmpl.Execute(w, &map[string]any{
		"ShowVersions": showVersions,
		"PackageIndex": &idx,
	})
	checkTmplError(err, "pkgInstallModalVersionsTmpl")
}

func (s *server) uninstall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := r.FormValue("packageName")
	var pkg v1alpha1.Package
	if err := s.pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
		return
	}
	// TODO: this should be changed to also broadcast the pending update first
	if err := uninstall.NewUninstaller(s.pkgClient).
		WithStatusWriter(statuswriter.Stderr()).
		Uninstall(ctx, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
	}
	s.broadcastPkgStatusUpdate(pkgName, nil, nil)
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
	err = pkgsPageTmpl.Execute(w, packages)
	checkTmplError(err, "packages")
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkgName"]
	pkg, status, manifest, err := describe.DescribePackage(r.Context(), pkgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred fetching package details of %v: \n%v\n", pkgName, err)
		return
	}
	err = pkgPageTmpl.Execute(w, &map[string]any{
		"Package":  pkg,
		"Status":   status,
		"Manifest": manifest,
	})
	checkTmplError(err, fmt.Sprintf("package-detail (%s)", pkgName))
}

func (s *server) supportPage(w http.ResponseWriter, r *http.Request) {
	if err := s.checkBootstrapped(r.Context()); err != nil {
		if err.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
			return
		}
		err := supportPageTmpl.Execute(w, err)
		checkTmplError(err, "support")
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *server) bootstrapPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == "POST" {
		client := bootstrap.NewBootstrapClient(
			s.restConfig,
			"",
			config.Version,
			bootstrap.BootstrapTypeAio,
		)
		if err := client.Bootstrap(ctx); err != nil {
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
			"Err":            err,
			"CurrentContext": s.rawConfig.CurrentContext,
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
		"ConfigErr":                 configErr,
		"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
		"CurrentContext":            currentContext,
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
	return nil
}

func (s *server) enrichContext(h http.Handler) http.Handler {
	return &handler.ContextEnrichingHandler{Source: s, Handler: h}
}

func (s *server) requireBootstrapped(h http.HandlerFunc) http.Handler {
	return &handler.PreconditionHandler{
		Precondition:  func(r *http.Request) error { return s.checkBootstrapped(r.Context()) },
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

func isPortConflictError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if osErr, ok := opErr.Err.(*os.SyscallError); ok {
			return osErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}
