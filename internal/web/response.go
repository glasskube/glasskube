package web

import (
	"encoding/json"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"net/http"
)

type Response struct {
	statusCode int
	templateId string
}

func (s *server) newToastResponse() *toast.ResponseBuilder {
	return toast.NewResponseBuilder(s.templates.toastTmpl).WithDismissible(true)
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
