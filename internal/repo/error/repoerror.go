package error

import (
	"errors"
)

type PartialErr struct {
	err error
}

func Partial(err error) *PartialErr {
	if err == nil {
		return nil
	}
	return &PartialErr{err: err}
}

func (e *PartialErr) Error() string {
	return e.err.Error()
}

func IsPartial(err error) bool {
	var e *PartialErr
	return err == nil || errors.As(err, &e)
}
