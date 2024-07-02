package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func (s *server) respondSuccess(w http.ResponseWriter) {
	err := s.templates.alertTmpl.Execute(w, map[string]any{
		"Message":     "Configuration updated successfully",
		"Dismissible": true,
		"Type":        "success",
	})
	checkTmplError(err, "success")
}

func (s *server) respondAlertAndLog(w http.ResponseWriter, err error, wrappingMsg string, alertType string) {
	if wrappingMsg != "" {
		err = fmt.Errorf("%v: %w", wrappingMsg, err)
	}
	fmt.Fprintf(os.Stderr, "%v\n", err)
	s.respondAlert(w, err.Error(), alertType)
}

func (s *server) respondAlert(w http.ResponseWriter, message string, alertType string) {
	w.Header().Add("Hx-Reselect", "div.alert") // overwrite any existing hx-select (which was a little intransparent sometimes)
	w.Header().Add("Hx-Reswap", "afterbegin")
	w.WriteHeader(http.StatusBadRequest)
	err := s.templates.alertTmpl.Execute(w, map[string]any{
		"Message":     message,
		"Dismissible": true,
		"Type":        alertType,
	})
	checkTmplError(err, "alert")
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
