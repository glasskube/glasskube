package web

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/glasskube/glasskube/internal/repo"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/glasskube/glasskube/pkg/kubeconfig"

	"github.com/glasskube/glasskube/pkg/manifest"

	"github.com/glasskube/glasskube/pkg/describe"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
	"github.com/gorilla/mux"

	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/open"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/uninstall"
)

//go:embed root
//go:embed templates
var embededFs embed.FS

type ServerOptions struct {
	Host       string
	Port       int32
	Kubeconfig string
}

type server struct {
	ServerOptions
	listener   net.Listener
	restConfig *rest.Config
	rawConfig  *api.Config
	pkgClient  *client.PackageV1Alpha1Client
	wsHub      *WsHub
	forwarders map[string]*open.OpenResult
	loadConfig func() (*rest.Config, *api.Config, error)
}

func NewServer(options ServerOptions) *server {
	server := server{
		ServerOptions: options,
		forwarders:    make(map[string]*open.OpenResult),
	}
	server.loadFileConfig()
	return &server
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

func (s *server) Start(ctx context.Context) error {
	if s.listener != nil {
		return errors.New("server is already listening")
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
	router.HandleFunc("/bootstrap", s.bootstrapPage)
	router.HandleFunc("/packages", s.packages)
	router.HandleFunc("/packages/install", s.install)
	router.HandleFunc("/packages/install/modal", s.installModal)
	router.HandleFunc("/packages/install/modal/versions", s.installModalVersions)
	router.HandleFunc("/packages/uninstall", s.uninstall)
	router.HandleFunc("/packages/open", s.open)
	router.HandleFunc("/packages/{pkgName}", s.packageDetail)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/packages", http.StatusFound)
	})
	http.Handle("/", router)

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
	ctx := s.prepareContext(r.Context())
	selectedVersion := r.FormValue("selectedVersion")
	go func() {
		status, err := install.NewInstaller(s.pkgClient).
			WithStatusWriter(statuswriter.Stderr()).
			InstallBlocking(ctx, pkgName, selectedVersion)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred installing %v: \n%v\n", pkgName, err)
			return
		}
		manifest, err := manifest.GetInstalledManifest(ctx, pkgName)
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

	_, _, manifest, err := describe.DescribePackage(s.prepareContext(r.Context()), pkgName)
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
	ctx := s.prepareContext(r.Context())
	pkgName := r.FormValue("packageName")
	var pkg v1alpha1.Package
	if err := s.pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
		return
	}
	// once we have blocking uninstall available, this should be changed to also broadcast the pending update first
	if err := uninstall.Uninstall(s.pkgClient, ctx, &pkg); err != nil {
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

	result, err := open.NewOpener().Open(s.prepareContext(r.Context()), pkgName, "")
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
	if err := s.checkPreconditions(r.Context()); err != nil {
		http.Redirect(w, r, "/support", http.StatusFound)
		return
	}

	packages, err := list.GetPackagesWithStatus(s.pkgClient, s.prepareContext(r.Context()), list.IncludePackageInfos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load packages: %v\n", err)
		return
	}
	err = pkgsPageTmpl.Execute(w, packages)
	checkTmplError(err, "packages")
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	if err := s.checkPreconditions(r.Context()); err != nil {
		http.Redirect(w, r, "/support", http.StatusFound)
		return
	}
	pkgName := mux.Vars(r)["pkgName"]
	pkg, status, manifest, err := describe.DescribePackage(s.prepareContext(r.Context()), pkgName)
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
	if err := s.checkPreconditions(r.Context()); err != nil {
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
	ctx := s.prepareContext(r.Context())
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
		if configErr := s.checkPreconditions(ctx); configErr != nil {
			// try to validate bootstrap again, if it failed last time:
			if configErr.BootstrapMissing() {
				isBootstrapped, bootstrapErr := bootstrap.IsBootstrapped(ctx, s.restConfig)
				if bootstrapErr != nil {
					fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", bootstrapErr)
				} else if isBootstrapped {
					configErr = nil
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}
			tplErr := bootstrapPageTmpl.Execute(w, &map[string]any{
				"ConfigErr":      configErr,
				"CurrentContext": s.rawConfig.CurrentContext,
			})
			checkTmplError(tplErr, "bootstrap")
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func (server *server) loadFileConfig() {
	server.loadConfig = func() (*rest.Config, *api.Config, error) { return kubeconfig.New(server.Kubeconfig) }
}

func (server *server) checkPreconditions(ctx context.Context) ServerConfigError {
	if server.pkgClient == nil {
		if err := server.initKubeConfig(); err != nil {
			return err
		}
	}

	isBootstrapped, err := bootstrap.IsBootstrapped(ctx, server.restConfig)
	if !isBootstrapped || err != nil {
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
		}
		return bootstrapError(err)
	}

	return nil
}

func (server *server) initKubeConfig() ServerConfigError {
	restConfig, rawConfig, err := server.loadConfig()
	if err != nil {
		return kubeconfigError(err)
	}
	client, err := client.New(restConfig)
	if err != nil {
		return kubeconfigError(err)
	}
	server.restConfig = restConfig
	server.rawConfig = rawConfig
	server.pkgClient = client
	return nil
}

func (server *server) prepareContext(ctx context.Context) context.Context {
	return client.SetupContextWithClient(ctx, server.restConfig, server.rawConfig, server.pkgClient)
}

func isPortConflictError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if osErr, ok := opErr.Err.(*os.SyscallError); ok {
			return osErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}
