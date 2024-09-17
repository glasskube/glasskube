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
	return errors.As(err, &e)
}
func IsPartialOrNil(err error) bool {
	return err == nil || IsPartial(err)
}

func IsComplete(err error) bool {
	return !IsPartialOrNil(err)
}
