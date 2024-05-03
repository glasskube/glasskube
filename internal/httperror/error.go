package httperror

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

type statusError struct {
	status string
	code   int
}

var (
	StatusError = errors.New("wrong status code")
)

func (err statusError) Error() string {
	return fmt.Sprintf("%v: %v", StatusError.Error(), err.status)
}

func (err statusError) Unwrap() error {
	return StatusError
}

func IsErrorResponse(respose *http.Response) bool {
	return respose.StatusCode >= 400
}

func CheckResponse(response *http.Response) error {
	if IsErrorResponse(response) {
		return &statusError{response.Status, response.StatusCode}
	} else {
		return nil
	}
}

func Is(err error, code int) bool {
	var statusErr statusError
	if errors.As(err, &statusErr) {
		return statusErr.code == code

	}
	return false
}

func IsNotFound(err error) bool {
	return Is(err, http.StatusNotFound)
}

func IsNetworkError(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}
