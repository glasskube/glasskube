package toast

import (
	"net/http"
)

type Response struct {
	ToastInput
	StatusCode int
	Err        error
}

type ResponseOption func(options *Response)

func WithStatusCode(statusCode int) ResponseOption {
	return func(options *Response) {
		options.StatusCode = statusCode
	}
}

func WithCssClass(cssClass string) ResponseOption {
	return func(options *Response) {
		options.CssClass = cssClass
	}
}

func WithErr(err error) ResponseOption {
	return func(options *Response) {
		options.Err = err
	}
}

func WithMessage(message string) ResponseOption {
	return func(options *Response) {
		options.Message = message
	}
}

func WithDismissible(dismissible bool) ResponseOption {
	return func(options *Response) {
		options.Dismissible = dismissible
	}
}

// Apply sets some reasonable defaults: if only an error is given, status code will be 500 and css class will be danger.
// If no error is given, status 200 OK and the success class are assumed, and the given Message will be used.
// However, all parts (Message, status code, class) can be set individually too.
func (r *Response) Apply() {
	if r.Err != nil {
		if r.Message == "" {
			r.Message = r.Err.Error()
		}
		if r.StatusCode == 0 {
			r.StatusCode = http.StatusInternalServerError
		}
	}

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}

	if r.CssClass == "" {
		if r.StatusCode < 200 || r.StatusCode >= 300 {
			r.CssClass = "danger"
		} else {
			r.CssClass = "success"
		}
	}
}
