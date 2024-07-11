package toast

import (
	"fmt"
	"github.com/glasskube/glasskube/internal/web/util"
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

type ResponseBuilder struct {
	toastResponse
	template *template.Template
}

func NewResponseBuilder(template *template.Template) *ResponseBuilder {
	return &ResponseBuilder{template: template}
}

func (b *ResponseBuilder) WithStatus(statusCode int) *ResponseBuilder {
	b.statusCode = statusCode
	return b
}

func (b *ResponseBuilder) WithDisplayClass(cssClass string) *ResponseBuilder {
	b.cssClass = cssClass
	return b
}

func (b *ResponseBuilder) WithErr(err error) *ResponseBuilder {
	b.err = err
	return b
}

func (b *ResponseBuilder) WithMessage(message string) *ResponseBuilder {
	b.message = message
	return b
}

func (b *ResponseBuilder) WithDismissible(dismissible bool) *ResponseBuilder {
	b.dismissible = dismissible
	return b
}

// Send sends the rendered toast to the given response writer, applying some reasonable defaults: if only an error
// is given, status code will be 500 and css class will be danger. If no error is given, status 200 OK and the
// success class are assumed, and the given message will be used. However, all parts (message, status code, class) can
// be set individually too.
func (b *ResponseBuilder) Send(w http.ResponseWriter) {
	// htmx headers to overwrite any existing/inherited hx-select, hx-swap, hx-target on the client
	w.Header().Add("Hx-Reselect", "div.toast")
	w.Header().Add("Hx-Reswap", "afterbegin")
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
		"cssClass":    b.cssClass,
	})
	util.CheckTmplError(err, "toast")

	if b.err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", b.err)
	}
}
