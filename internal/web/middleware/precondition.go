package middleware

import "net/http"

type Precondition func(r *http.Request) error
type PreconditionFailedHandleFunc func(w http.ResponseWriter, r *http.Request, err error)

type PreconditionHandler struct {
	Precondition  Precondition
	Handler       http.Handler
	FailedHandler PreconditionFailedHandleFunc
}

func (ph *PreconditionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ph.Precondition(r); err != nil {
		ph.FailedHandler(w, r, err)
	} else {
		ph.Handler.ServeHTTP(w, r)
	}
}
