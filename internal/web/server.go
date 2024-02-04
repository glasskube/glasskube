package web

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/glasskube/glasskube/internal/web/components"

	"github.com/glasskube/glasskube/internal/cliutils"

	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/list"
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

	Host = "localhost"
	Port = 8580
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

func Start(ctx context.Context, support *ServerConfigSupport) error {
	root, err := fs.Sub(embededFs, "root")
	if err != nil {
		return err
	}

	wsHub := NewHub()

	fileServer := http.FileServer(http.FS(root))
	http.Handle("/static/", fileServer)
	http.Handle("/favicon.ico", fileServer)
	http.HandleFunc("/ws", wsHub.handler)
	http.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		pkgClient := client.FromContext(ctx)
		pkgName := r.FormValue("packageName")
		go func() {
			status, err := install.NewInstaller(pkgClient).
				WithStatusWriter(statuswriter.Stderr()).
				InstallBlocking(ctx, pkgName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred installing %v: \n%v\n", pkgName, err)
			}

			// broadcast the status update to all clients
			var bf bytes.Buffer
			components.RenderInstallButton(&bf, installBtnTmpl, pkgName, status)
			wsHub.Broadcast <- bf.Bytes()
		}()

		// broadcast the pending button to all clients (note that we do not return any html from the install endpoint)
		var bf bytes.Buffer
		components.RenderInstallButton(&bf, installBtnTmpl, pkgName, &client.PackageStatus{
			Status: "Pending",
		})
		wsHub.Broadcast <- bf.Bytes()
	})
	http.HandleFunc("/uninstall", func(w http.ResponseWriter, r *http.Request) {
		pkgClient := client.FromContext(ctx)
		pkgName := r.FormValue("packageName")
		pkg, err := list.Get(pkgClient, ctx, pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
			return
		}

		err = uninstall.Uninstall(pkgClient, ctx, pkg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
		}

		// broadcast the button depending on status to all clients
		var bf bytes.Buffer
		components.RenderInstallButton(&bf, installBtnTmpl, pkgName, nil)
		wsHub.Broadcast <- bf.Bytes()
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if support != nil {
			err := supportPageTmpl.Execute(w, support)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred rendering the response: \n%v\n", err)
			}
			return
		}

		pkgClient := client.FromContext(ctx)
		packages, _ := list.GetPackagesWithStatus(pkgClient, ctx, false)
		err := pkgsPageTmpl.Execute(w, packages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred rendering the response: \n%v\n", err)
		}
	})

	bindAddr := fmt.Sprintf("%v:%d", Host, Port)
	var listener net.Listener

	listener, err = net.Listen("tcp", bindAddr)
	if err != nil {
		// Checks if Port Conflict Error exists
		if isPortConflictError(err) {
			userInput := cliutils.YesNoPrompt(
				"Port is already in use.\nShould glasskube use a different port? (Y/n): ", true)
			if userInput {
				listener, err = net.Listen("tcp", ":0")
				if err != nil {
					panic(err)
				}
				bindAddr = fmt.Sprintf("%v:%d", Host, listener.Addr().(*net.TCPAddr).Port)
			} else {
				fmt.Println("Exiting. User chose not to use a different port.")
				os.Exit(1)
			}
		} else {
			// If no Port Conflict error is found, return other errors
			return err
		}
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing listener: %v\n", err)
		}
	}()

	fmt.Printf("glasskube UI is available at http://%v\n", bindAddr)
	_ = openInBrowser("http://" + bindAddr)

	go wsHub.Run()
	srv := &http.Server{}
	err = srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func openInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func isPortConflictError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if osErr, ok := opErr.Err.(*os.SyscallError); ok {
			return osErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}
