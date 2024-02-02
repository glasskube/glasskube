package web

import (
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

	"github.com/glasskube/glasskube/internal/cliutils"

	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/uninstall"
	"k8s.io/apimachinery/pkg/api/errors"
)

var Host = "localhost"
var Port = 8580

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

func Start(ctx context.Context, support *ServerConfigSupport) error {
	pkgTemplate, err := template.ParseFS(embededFs, "templates/packages.html")
	if err != nil {
		return err
	}
	supportTemplate, err := template.ParseFS(embededFs, "templates/support.html")
	if err != nil {
		return err
	}

	root, err := fs.Sub(embededFs, "root")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(root))
	http.Handle("/static/", fileServer)
	http.Handle("/favicon.ico", fileServer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if support != nil {
			err := supportTemplate.Execute(w, support)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred rendering the response: \n%v\n", err)
			}
			return
		}

		pkgClient := client.FromContext(ctx)
		if r.Method == "POST" {
			pkgName := r.FormValue("packageName")
			pkg, err := list.Get(pkgClient, ctx, pkgName)
			if err != nil && !errors.IsNotFound(err) {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}
			if pkg != nil {
				err := uninstall.Uninstall(pkgClient, ctx, pkg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "An error occurred uninstalling %v: \n%v\n", pkgName, err)
				}
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				err := install.NewInstaller(pkgClient).
					WithStatusWriter(statuswriter.Stderr()).
					Install(ctx, pkgName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "An error occurred installing %v: \n%v\n", pkgName, err)
				}
				http.Redirect(w, r, "/", http.StatusFound)
			}
			return
		}

		packages, _ := list.GetPackagesWithStatus(pkgClient, ctx, false)
		err := pkgTemplate.Execute(w, packages)
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
