package web

import (
	"context"
	"embed"
	"fmt"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/uninstall"
	"html/template"
	"io/fs"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
				fmt.Fprintf(os.Stderr, "An error occured rendering the response: \n%v\n", err)
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
					fmt.Fprintf(os.Stderr, "An error occured uninstalling %v: \n%v\n", pkgName, err)
				}
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				_, err := install.Install(pkgClient, ctx, pkgName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "An error occured installing %v: \n%v\n", pkgName, err)
				}
				http.Redirect(w, r, "/", http.StatusFound)
			}
			return
		}

		packages, _ := list.GetPackagesWithStatus(pkgClient, ctx, false)
		err := pkgTemplate.Execute(w, packages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured rendering the response: \n%v\n", err)
		}
	})

	bindAddr := fmt.Sprintf("%v:%d", Host, Port)
	url := fmt.Sprintf("http://%v", bindAddr)
	fmt.Printf("glasskube UI is available at %v\n", url)
	_ = openInBrowser(url)

	err = http.ListenAndServe(bindAddr, nil)
	if err != nil {
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
