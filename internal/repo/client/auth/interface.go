package auth

import "net/http"

type Authenticator interface {
	Authenticate(request *http.Request)
}

func Basic(username, password string) Authenticator {
	return &basicAuthenticator{username: username, password: password}
}

func Bearer(token string) Authenticator {
	return &bearerAuthenticator{token: token}
}

func Noop() Authenticator {
	return NoopAuthenticator{}
}
