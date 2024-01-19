package web

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

var Host = "localhost"
var Port = 8580

func Start() error {
	bindAddr := fmt.Sprintf("%v:%d", Host, Port)
	url := fmt.Sprintf("http://%v", bindAddr)
	fmt.Printf("glasskube UI is available at %v\n", url)
	_ = openInBrowser(url)

	http.HandleFunc("/", index)
	err := http.ListenAndServe(bindAddr, nil)
	if err != nil {
		return err
	}
	return nil
}

func index(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World from Glasskube")
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
