package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type toastResponse struct {
	statusCode  int
	cssClass    string
	err         error
	message     string
	dismissible bool
}

type toastResponseBuilder struct {
	toastResponse
	template *template.Template
}

func (s *server) newToastResponse() *toastResponseBuilder {
	return &toastResponseBuilder{
		toastResponse{
			message:     "",
			dismissible: true,
		},
		s.templates.alertTmpl,
	}
}

func (b *toastResponseBuilder) withStatus(statusCode int) *toastResponseBuilder {
	b.statusCode = statusCode
	return b
}

func (b *toastResponseBuilder) withDisplayClass(cssClass string) *toastResponseBuilder {
	b.cssClass = cssClass
	return b
}

func (b *toastResponseBuilder) withErr(err error) *toastResponseBuilder {
	b.err = err
	return b
}

func (b *toastResponseBuilder) withMessage(message string) *toastResponseBuilder {
	b.message = message
	return b
}

func (b *toastResponseBuilder) withDismissible(dismissible bool) *toastResponseBuilder {
	b.dismissible = dismissible
	return b
}

func (b *toastResponseBuilder) send(w http.ResponseWriter) {
	w.Header().Add("Hx-Reselect", "div.toast") // overwrite any existing/inherited hx-select on the client
	w.Header().Add("Hx-Reswap", "afterbegin")  // overwrite any existing/inherited hx-swap on the client
	w.Header().Add("Hx-Retarget", "#toast-container")

	if b.err != nil {
		if b.message == "" {
			b.message = b.err.Error()
		}
		if b.statusCode == 0 {
			b.statusCode = http.StatusInternalServerError
		}
	}

	if b.statusCode == 0 {
		b.statusCode = http.StatusOK
	}
	w.WriteHeader(b.statusCode)

	if b.cssClass == "" {
		if b.statusCode < 200 || b.statusCode >= 300 {
			b.cssClass = "danger"
		} else {
			b.cssClass = "success"
		}
	}

	err := b.template.Execute(w, map[string]any{
		"Message":     b.message,
		"Dismissible": b.dismissible,
		"Type":        b.cssClass,
	})
	checkTmplError(err, "toast")

	if b.err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", b.err)
	}
}

func (s *server) respondWithToast(w http.ResponseWriter, toast *toastResponse) {
	w.Header().Add("Hx-Reselect", "div.toast") // overwrite any existing/inherited hx-select on the client
	w.Header().Add("Hx-Reswap", "afterbegin")  // overwrite any existing/inherited hx-swap on the client
	w.WriteHeader(toast.statusCode)

	// TODO maybe this logic should go into the build() function?
	toastMessage := toast.message
	if toast.err != nil && toast.message == "" {
		toastMessage = toast.err.Error()
		// TODO maybe set class to danger ?? not sure yet
	}

	class := "success"
	if toast.cssClass != "" {
		class = toast.cssClass
	} else if toast.statusCode < 200 || toast.statusCode >= 300 {
		class = "danger"
	}

	err := s.templates.alertTmpl.Execute(w, map[string]any{
		"Message":     toastMessage,
		"Dismissible": toast.dismissible,
		"Type":        class,
	})
	checkTmplError(err, "toast")

	if toast.err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", toast.err)
	}
}

func (s *server) respondErrorToastAndLog(w http.ResponseWriter, err error, wrappingMsg string, alertType string) {
	if wrappingMsg != "" {
		err = fmt.Errorf("%v: %w", wrappingMsg, err)
	}
	fmt.Fprintf(os.Stderr, "%v\n", err)
	s.respondToast(w, err.Error(), alertType)
}

func (s *server) respondToast(w http.ResponseWriter, message string, toastType string) {
	w.Header().Add("Hx-Reselect", "div.toast") // overwrite any existing hx-select (which was a little intransparent sometimes)
	w.Header().Add("Hx-Reswap", "afterbegin")
	w.WriteHeader(http.StatusBadRequest)
	err := s.templates.alertTmpl.Execute(w, map[string]any{
		"Message":     message,
		"Dismissible": true,
		"Type":        toastType,
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
