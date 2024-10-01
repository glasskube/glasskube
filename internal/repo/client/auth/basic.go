package auth

import (
	"net/http"
)

type basicAuthenticator struct {
	username, password string
}

// Authenticate implements client.Authenticator.
func (b *basicAuthenticator) Authenticate(request *http.Request) {
	request.SetBasicAuth(b.username, b.password)
}
