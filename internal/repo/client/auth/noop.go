package auth

import (
	"net/http"
)

type NoopAuthenticator struct{}

// Authenticate implements client.Authenticator.
func (n NoopAuthenticator) Authenticate(request *http.Request) {
}
