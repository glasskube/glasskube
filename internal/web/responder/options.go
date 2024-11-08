package responder

import "github.com/glasskube/glasskube/internal/web/types"

type response struct {
	statusCode   int
	templateData any
	partialErr   error
}

// TODO merge with toast response option???
type ResponseOption func(*response)

func WithStatusCode(statusCode int) ResponseOption {
	return func(r *response) {
		r.statusCode = statusCode
	}
}

func ContextualizedTemplate(templateData types.ContextInjectable) ResponseOption {
	return func(r *response) {
		r.templateData = templateData
	}
}

func RawTemplate(templateData any) ResponseOption {
	return func(r *response) {
		r.templateData = templateData
	}
}

func WithPartialErr(partialErr error) ResponseOption {
	return func(r *response) {
		r.partialErr = partialErr
	}
}
