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

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/bootstrap"

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

type ServerConfigSupport struct {
	KubeconfigMissing         bool
	KubeconfigDefaultLocation string
	KubeconfigError           error
	BootstrapMissing          bool
	BootstrapCheckError       error
}

func init() {
	loadTemplates()
}

type server struct {
	host       string
	port       int32
	listener   net.Listener
	ctx        context.Context
	cfg        *rest.Config
	rawCfg     *api.Config
	support    *ServerConfigSupport
	pkgClient  *client.PackageV1Alpha1Client
	wsHub      *WsHub
	forwarders map[string]*open.OpenResult
}

func NewServer(host string, port int32) *server {
	return &server{
		forwarders: make(map[string]*open.OpenResult),
		host:       host,
		port:       port,
	}
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

func (s *server) Start(ctx context.Context, support *ServerConfigSupport) error {
	if s.listener != nil {
		return errors.New("server is already listening")
	}

	s.cfg = client.ConfigFromContext(ctx)
	s.rawCfg = client.RawConfigFromContext(ctx)

	if support == nil {
		isBootstrapped, err := bootstrap.IsBootstrapped(ctx, s.cfg)
		if !isBootstrapped || err != nil {
			support = &ServerConfigSupport{
				BootstrapMissing:    !isBootstrapped,
				BootstrapCheckError: err,
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
		}
	}

	s.support = support
	s.ctx = ctx
	s.pkgClient = client.FromContext(ctx)

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
	router.HandleFunc("/packages/uninstall", s.uninstall)
	router.HandleFunc("/packages/open", s.open)
	router.HandleFunc("/packages/{pkgName}", s.packageDetail)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/packages", http.StatusFound)
	})
	http.Handle("/", router)

	bindAddr := fmt.Sprintf("%v:%d", s.host, s.port)

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
				bindAddr = fmt.Sprintf("%v:%d", s.host, s.listener.Addr().(*net.TCPAddr).Port)
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
	go func() {
		// TODO: Add version
		status, err := install.NewInstaller(s.pkgClient).
			WithStatusWriter(statuswriter.Stderr()).
			InstallBlocking(s.ctx, pkgName, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred installing %v: \n%v\n", pkgName, err)
			return
		}
		manifest, err := manifest.GetInstalledManifest(s.ctx, pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not fetch manifest of %v: %v\n", pkgName, err)
			return
		}
		s.broadcastPkgStatusUpdate(pkgName, status, manifest)
	}()
	s.broadcastPkgStatusUpdate(pkgName, client.NewPendingStatus(), nil)
}

func (s *server) uninstall(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	pkg, err := list.Get(s.pkgClient, s.ctx, pkgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
		return
	}
	// once we have blocking uninstall available, this should be changed to also broadcast the pending update first
	err = uninstall.Uninstall(s.pkgClient, s.ctx, pkg)
	if err != nil {
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

	result, err := open.NewOpener().Open(s.ctx, pkgName, "")
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
	if s.support != nil {
		http.Redirect(w, r, "/support", http.StatusFound)
		return
	}

	packages, err := list.GetPackagesWithStatus(s.pkgClient, s.ctx, list.IncludePackageInfos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load packages: %v\n", err)
		return
	}
	err = pkgsPageTmpl.Execute(w, packages)
	checkTmplError(err, "packages")
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	if s.support != nil {
		http.Redirect(w, r, "/support", http.StatusFound)
		return
	}
	pkgName := mux.Vars(r)["pkgName"]
	pkg, status, manifest, err := describe.DescribePackage(s.ctx, pkgName)
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
	if s.support != nil {
		if s.support.BootstrapMissing {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
			return
		}
		err := supportPageTmpl.Execute(w, s.support)
		checkTmplError(err, "support")
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *server) bootstrapPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		client := bootstrap.NewBootstrapClient(
			s.cfg,
			"",
			config.Version,
			bootstrap.BootstrapTypeAio,
		)
		if err := client.Bootstrap(s.ctx); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			err := bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-failure", nil)
			checkTmplError(err, "bootstrap-failure")
		} else {
			s.support = nil
			err := bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-success", nil)
			checkTmplError(err, "bootstrap-success")
		}
	} else {
		if s.support != nil && s.support.BootstrapMissing {
			// try to validate bootstrap again, if it failed last time:
			if s.support.BootstrapCheckError != nil {
				isBootstrapped, err := bootstrap.IsBootstrapped(s.ctx, s.cfg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
				} else if isBootstrapped {
					s.support = nil
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}
			err := bootstrapPageTmpl.Execute(w, &map[string]any{
				"Support":        s.support,
				"CurrentContext": s.rawCfg.CurrentContext,
			})
			checkTmplError(err, "bootstrap")
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
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
