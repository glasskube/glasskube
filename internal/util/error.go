package util

import "errors"

func GetRootCause(err error) error {
	unwrapped := errors.Unwrap(err)
	if unwrapped != nil {
		return GetRootCause(unwrapped)
	} else {
		return err
	}
}
