package web

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/web/components"

	"github.com/glasskube/glasskube/internal/cliutils"

	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/open"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/uninstall"
)

var (
	baseTemplate    *template.Template
	pkgsPageTmpl    *template.Template
	supportPageTmpl *template.Template
	installBtnTmpl  *template.Template

	//go:embed root
	//go:embed templates
	embededFs embed.FS
)

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

func loadTemplates() {
	templateFuncs := template.FuncMap{
		"ToInstallButtonInput": components.ToInstallButtonInput,
	}
	baseTemplate = template.Must(
		template.New("base.html").Funcs(templateFuncs).ParseFS(embededFs, "templates/layout/base.html"))
	pkgsPageTmpl = template.Must(template.Must(baseTemplate.Clone()).
		ParseFS(embededFs, "templates/pages/packages.html", "templates/components/*.html"))
	supportPageTmpl = template.Must(template.Must(baseTemplate.Clone()).
		ParseFS(embededFs, "templates/pages/support.html", "templates/components/*.html"))
	installBtnTmpl = template.Must(template.New("installbutton").
		ParseFS(embededFs, "templates/components/installbutton.html"))
}

type server struct {
	host       string
	port       int32
	listener   net.Listener
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

func (s *server) Start(ctx context.Context, support *ServerConfigSupport) error {
	if s.listener != nil {
		return errors.New("server is already listening")
	}

	s.pkgClient = client.FromContext(ctx)

	root, err := fs.Sub(embededFs, "root")
	if err != nil {
		return err
	}

	s.wsHub = NewHub()

	fileServer := http.FileServer(http.FS(root))
	http.Handle("/static/", fileServer)
	http.Handle("/favicon.ico", fileServer)
	http.HandleFunc("/ws", s.wsHub.handler)
	http.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		pkgName := r.FormValue("packageName")
		go func() {
			status, err := install.NewInstaller(s.pkgClient).
				WithStatusWriter(statuswriter.Stderr()).
				InstallBlocking(ctx, pkgName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred installing %v: \n%v\n", pkgName, err)
			}

			var packageInfo v1alpha1.PackageInfo
			var manifest *v1alpha1.PackageManifest
			// TODO: Change this to use the actual package info name instead of the package name
			if err := s.pkgClient.PackageInfos().Get(ctx, pkgName, &packageInfo); err != nil {
				fmt.Printf("could not fetch PackageInfo %v: %v\n", pkgName, err)
			} else {
				manifest = packageInfo.Status.Manifest
			}

			// broadcast the status update to all clients
			var bf bytes.Buffer
			components.RenderInstallButton(&bf, installBtnTmpl, pkgName, status, manifest)
			s.wsHub.Broadcast <- bf.Bytes()
		}()

		// broadcast the pending button to all clients (note that we do not return any html from the install endpoint)
		var bf bytes.Buffer
		components.RenderInstallButton(&bf, installBtnTmpl, pkgName, &client.PackageStatus{Status: "Pending"}, nil)
		s.wsHub.Broadcast <- bf.Bytes()
	})
	http.HandleFunc("/uninstall", func(w http.ResponseWriter, r *http.Request) {
		pkgName := r.FormValue("packageName")
		pkg, err := list.Get(s.pkgClient, ctx, pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
			return
		}

		err = uninstall.Uninstall(s.pkgClient, ctx, pkg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
		}

		// broadcast the button depending on status to all clients
		var bf bytes.Buffer
		components.RenderInstallButton(&bf, installBtnTmpl, pkgName, nil, nil)
		s.wsHub.Broadcast <- bf.Bytes()
	})
	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		pkgName := r.FormValue("packageName")
		if result, ok := s.forwarders[pkgName]; ok {
			result.WaitReady()
			_ = cliutils.OpenInBrowser(result.Url)
			return
		}

		result, err := open.NewOpener().Open(ctx, pkgName, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open %v: %v\n", pkgName, err)
		} else {
			s.forwarders[pkgName] = result
			result.WaitReady()
			_ = cliutils.OpenInBrowser(result.Url)
			w.WriteHeader(http.StatusAccepted)
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if support != nil {
			err := supportPageTmpl.Execute(w, support)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred rendering the response: \n%v\n", err)
			}
			return
		}

		packages, _ := list.GetPackagesWithStatus(s.pkgClient, ctx, list.IncludePackageInfos)
		err := pkgsPageTmpl.Execute(w, packages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred rendering the response: \n%v\n", err)
		}
	})

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

func isPortConflictError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if osErr, ok := opErr.Err.(*os.SyscallError); ok {
			return osErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}
