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

func IsErrorResponse(response *http.Response) bool {
	return response.StatusCode >= 400
}

func CheckResponse(response *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return response, err
	} else if IsErrorResponse(response) {
		return response, &statusError{response.Status, response.StatusCode}
	} else {
		return response, nil
	}
}

func Is(err error, code int) bool {
	var statusErr *statusError
	if errors.As(err, &statusErr) {
		return statusErr.code == code
	}
	return false
}

func IsNotFound(err error) bool {
	return Is(err, http.StatusNotFound)
}

func IsTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
