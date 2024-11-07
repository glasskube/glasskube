package responder

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

func WithTemplateData(templateData any) ResponseOption {
	return func(r *response) {
		r.templateData = templateData
	}
}

func WithPartialErr(partialErr error) ResponseOption {
	return func(r *response) {
		r.partialErr = partialErr
	}
}

func (r *response) Apply() {

}
