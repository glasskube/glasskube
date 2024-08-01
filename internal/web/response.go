package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/util"
)

// sendToast builds a toast from the given options and sends it to the given response writer. If the response
// contains an error, this is also logged to stderr.
func (s *server) sendToast(w http.ResponseWriter, options ...toast.ResponseOption) {
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

	err := s.templates.toastTmpl.Execute(w, response.ToastInput)
	util.CheckTmplError(err, "toast")

	if response.Err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", response.Err)
	}
}

// swappingRedirect adds the Hx-Location header to the response, which, when interpreted by htmx.js, will make
// the frontend redirect to the given path and swap the given target with the given slect from the response
// also see: https://htmx.org/headers/hx-location/
func (s *server) swappingRedirect(w http.ResponseWriter, path string, target string, slect string) {
	locationData := map[string]string{
		"path":   path,
		"target": target,
		"select": slect,
	}
	locationJson, _ := json.Marshal(locationData)
	w.Header().Add("Hx-Location", string(locationJson))
}

func (s *server) sendYamlModal(w http.ResponseWriter, obj string, alertContent any) {
	// htmx headers to overwrite any existing/inherited hx-select, hx-swap, hx-target on the client
	w.Header().Add("Hx-Reselect", "#yaml-modal")
	w.Header().Add("Hx-Reswap", "innerHTML")
	w.Header().Add("Hx-Retarget", "#modal-container")

	w.WriteHeader(http.StatusOK)

	e := s.templates.yamlModalTmpl.Execute(w, map[string]any{
		"AlertContent": alertContent,
		"Object":       obj,
	})
	util.CheckTmplError(e, "yaml-modal")
}
