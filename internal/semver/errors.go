package semver

import "go.uber.org/multierr"

type ConstraintValidationError []error

func ErrConstraintValidation(errors []error) error {
	if len(errors) > 0 {
		return (*ConstraintValidationError)(&errors)
	} else {
		return nil
	}
}

func (err *ConstraintValidationError) Error() string {
	return multierr.Combine(*err...).Error()
}

func (err *ConstraintValidationError) Unwrap() []error {
	return *err
}

func (err *ConstraintValidationError) Is(other error) bool {
	_, ok := other.(*ConstraintValidationError)
	return ok
}
