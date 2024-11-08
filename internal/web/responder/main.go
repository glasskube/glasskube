package responder

import (
	"encoding/json"
	"fmt"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/types"
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

type htmlResponder struct {
	contextProvider ContextProvider
	templates       *templates
	cloudId         string
}

var responder *htmlResponder

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

	responder = &htmlResponder{
		contextProvider: contextProvider,
		templates:       t,
		cloudId:         telemetry.GetMachineId(),
	}
}

func SendPage(w http.ResponseWriter, r *http.Request, templateName string, options ...ResponseOption) {
	responder.sendPage(w, r, templateName, options...)
}

func (res *htmlResponder) sendPage(w io.Writer, req *http.Request, templateName string, options ...ResponseOption) {
	r := &response{}
	for _, opt := range options {
		opt(r)
	}

	navbar := types.Navbar{}
	if pathParts := strings.Split(req.URL.Path, "/"); len(pathParts) >= 2 {
		navbar.ActiveItem = pathParts[1]
	}
	r.templateData.SetContextData(types.TemplateContextData{
		Navbar:             navbar,
		VersionDetails:     types.VersionDetails{}, // TODO from server (also think about caching when getting the version!!)
		CurrentContext:     res.contextProvider.GetCurrentContext(),
		GitopsMode:         res.contextProvider.IsGitopsModeEnabled(),
		Error:              r.partialErr,
		CacheBustingString: config.Version,
		CloudId:            res.cloudId,
		TemplateName:       templateName,
	})

	tmplErr := res.templates.baseTemplate.ExecuteTemplate(w, "base.html", r.templateData)
	// TODO tmpl error should return status 500 ??
	checkTmplError(tmplErr, templateName)
}

func SendComponent(w http.ResponseWriter, r *http.Request, templateName string, options ...ResponseOption) {
	responder.sendComponent(w, r, templateName, options...)
}

func (res *htmlResponder) sendComponent(w io.Writer, req *http.Request, templateName string, options ...ResponseOption) {
	r := &response{}
	for _, opt := range options {
		opt(r)
	}

	// TODO leave "reduced" TemplateContextData or simply send all of it every time ?
	r.templateData.SetContextData(types.TemplateContextData{
		CurrentContext: res.contextProvider.GetCurrentContext(),
		GitopsMode:     res.contextProvider.IsGitopsModeEnabled(),
		Error:          r.partialErr,
		CloudId:        res.cloudId,
		TemplateName:   templateName,
	})

	tmplErr := res.templates.baseTemplate.ExecuteTemplate(w, templateName, r.templateData)
	checkTmplError(tmplErr, templateName)
}

// SendToast builds a toast from the given options and sends it to the given response writer. If the response
// contains an error, this is also logged to stderr.
func SendToast(w http.ResponseWriter, options ...toast.ResponseOption) {
	responder.sendToast(w, options...)
}

func (res *htmlResponder) sendToast(w http.ResponseWriter, options ...toast.ResponseOption) {
	response := toast.Response{ToastInput: toast.ToastInput{Dismissible: true}}
	for _, opt := range options {
		opt(&response)
	}
	response.Apply()

	// htmx headers to overwrite any existing/inherited hx-select, hx-swap, hx-target on the client
	w.Header().Add(hxReselect, "div.toast")
	w.Header().Add(hxReswap, "afterbegin")
	w.Header().Add(hxRetarget, "#toast-container")

	w.WriteHeader(response.StatusCode)

	tmplName := "components/toast"
	tmplErr := res.templates.baseTemplate.ExecuteTemplate(w, tmplName, response.ToastInput)
	checkTmplError(tmplErr, tmplName)

	if response.Err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", response.Err)
	}
}

// Redirect adds the Hx-Location header to the response, which, when interpreted by htmx.js, will make
// the frontend redirect to the given path and swap the given target with the given slect from the response
// also see: https://htmx.org/headers/hx-location/
func Redirect(w http.ResponseWriter, path string) {
	responder.redirect(w, path)
}

func (res *htmlResponder) redirect(w http.ResponseWriter, path string) {
	target := "main"
	slect := "main"
	locationData := map[string]string{
		"path":   path,
		"target": target,
		"select": slect,
	}
	locationJson, _ := json.Marshal(locationData)
	w.Header().Add(hxLocation, string(locationJson))
}

func SendYamlModal(w http.ResponseWriter, obj string, alertContent any) {
	responder.sendYamlModal(w, obj, alertContent)
}

func (res *htmlResponder) sendYamlModal(w http.ResponseWriter, obj string, alertContent any) {
	// htmx headers to overwrite any existing/inherited hx-select, hx-swap, hx-target on the client
	w.Header().Add(hxReselect, "#yaml-modal")
	w.Header().Add(hxReswap, "innerHTML")
	w.Header().Add(hxRetarget, "#modal-container")

	w.WriteHeader(http.StatusOK)

	tmplName := "components/yaml-modal"
	tmplErr := res.templates.baseTemplate.ExecuteTemplate(w, tmplName, map[string]any{
		"AlertContent": alertContent,
		"Object":       obj,
	})
	checkTmplError(tmplErr, tmplName)
}

func checkTmplError(e error, tmplName string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "\nUnexpected error rendering '%v': %v\n – This is most likely a BUG – "+
			"Please report it here: https://github.com/glasskube/glasskube\n\n", tmplName, e)
	}
}
