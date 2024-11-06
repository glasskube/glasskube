package responder

import (
	"encoding/json"
	"fmt"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/util"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

type ContextProvider interface {
	GetCurrentContext() string
	IsGitopsModeEnabled() bool
}

type templateResponder struct {
	contextProvider ContextProvider
	templates       *templates
	cloudId         string
}

var responder *templateResponder

func Init(contextProvider ContextProvider, webFs fs.FS) {
	if responder != nil {
		panic("responder already initialized")
	}

	t := &templates{
		fs: webFs,
	}
	t.ParseTemplates()
	if config.IsDevBuild() {
		if err := t.WatchTemplates(); err != nil {
			fmt.Fprintf(os.Stderr, "templates will not be parsed after changes: %v\n", err)
		}
	}

	responder = &templateResponder{
		contextProvider: contextProvider,
		templates:       t,
		cloudId:         telemetry.GetMachineId(),
	}
}

func SendPage(w http.ResponseWriter, r *http.Request, templateName string, data any, err error) {
	// TODO options like in toast
	responder.sendPage(w, r, templateName, data, err)
}

func (res *templateResponder) sendPage(w io.Writer, r *http.Request, templateName string, data any, err error) {
	navbar := Navbar{}
	if pathParts := strings.Split(r.URL.Path, "/"); len(pathParts) >= 2 {
		navbar.ActiveItem = pathParts[1]
	}
	tmplErr := res.templates.baseTemplate.ExecuteTemplate(w, "base.html", Page{
		Navbar:             navbar,
		VersionDetails:     VersionDetails{}, // TODO from server
		CurrentContext:     res.contextProvider.GetCurrentContext(),
		GitopsMode:         res.contextProvider.IsGitopsModeEnabled(),
		Error:              err,
		CacheBustingString: config.Version,
		CloudId:            res.cloudId,
		TemplateName:       templateName,
		TemplateData:       data,
	})
	checkTmplError(tmplErr, templateName)
}

// sendToast builds a toast from the given options and sends it to the given response writer. If the response
// contains an error, this is also logged to stderr.
func SendToast(w http.ResponseWriter, options ...toast.ResponseOption) {
	responder.sendToast(w, options...)
}

func (res *templateResponder) sendToast(w http.ResponseWriter, options ...toast.ResponseOption) {
	response := toast.Response{ToastInput: toast.ToastInput{Dismissible: true}}
	for _, opt := range options {
		opt(&response)
	}
	response.Apply()

	// htmx headers to overwrite any existing/inherited hx-select, hx-swap, hx-target on the client
	w.Header().Add("Hx-Reselect", "div.toast")
	w.Header().Add("Hx-Reswap", "afterbegin")
	w.Header().Add("Hx-Retarget", "#toast-container")

	w.WriteHeader(response.StatusCode)

	tmplErr := res.templates.baseTemplate.ExecuteTemplate(w, "components/toast", response.ToastInput)
	util.CheckTmplError(tmplErr, "toast")

	if response.Err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", response.Err)
	}
}

// swappingRedirect adds the Hx-Location header to the response, which, when interpreted by htmx.js, will make
// the frontend redirect to the given path and swap the given target with the given slect from the response
// also see: https://htmx.org/headers/hx-location/
func (res *templateResponder) swappingRedirect(w http.ResponseWriter, path string, target string, slect string) {
	locationData := map[string]string{
		"path":   path,
		"target": target,
		"select": slect,
	}
	locationJson, _ := json.Marshal(locationData)
	w.Header().Add("Hx-Location", string(locationJson))
}

func (res *templateResponder) sendYamlModal(w http.ResponseWriter, obj string, alertContent any) {
	// htmx headers to overwrite any existing/inherited hx-select, hx-swap, hx-target on the client
	w.Header().Add("Hx-Reselect", "#yaml-modal")
	w.Header().Add("Hx-Reswap", "innerHTML")
	w.Header().Add("Hx-Retarget", "#modal-container")

	w.WriteHeader(http.StatusOK)

	e := res.templates.baseTemplate.ExecuteTemplate(w, "components/yaml-modal", map[string]any{
		"AlertContent": alertContent,
		"Object":       obj,
	})
	util.CheckTmplError(e, "yaml-modal")
}

func checkTmplError(e error, tmplName string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "\nUnexpected error rendering '%v': %v\n – This is most likely a BUG – "+
			"Please report it here: https://github.com/glasskube/glasskube\n\n", tmplName, e)
	}
}
