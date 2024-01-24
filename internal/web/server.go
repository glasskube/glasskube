package web

import (
	"context"
	"embed"
	"fmt"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
	"html/template"
	"io/fs"
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

func Start(ctx context.Context) error {
	tmpl, err := template.ParseFS(embededFs, "templates/packages.html")
	if err != nil {
		return err
	}

	root, _ := fs.Sub(embededFs, "root")
	fileServer := http.FileServer(http.FS(root))
	http.Handle("/static/", fileServer)
	http.Handle("/favicon.ico", fileServer)
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		// TODO check whether to show kubeconfig helper page + reload button #31

		pkgClient := client.FromContext(ctx)
		packages, _ := list.GetPackagesWithStatus(pkgClient, ctx, false)
		err := tmpl.Execute(w, packages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
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
