package auth

import (
	"net/http"
)

type bearerAuthenticator struct {
	token string
}

// Authenticate implements client.Authenticator.
func (b *bearerAuthenticator) Authenticate(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+b.token)
}
