package manifestvalues

import (
	"errors"
	"fmt"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

type KeyError struct {
	key string
}

func (err *KeyError) Error() string {
	return fmt.Sprintf("no such key: %v", err.key)
}

func NewKeyError(key string) error {
	return &KeyError{key: key}
}

type ValueTypeError struct {
	t v1alpha1.ValueType
}

func (err *ValueTypeError) Error() string {
	return fmt.Sprintf("unhandled type: %v (this is a bug)", err.t)
}

func NewValueTypeError(t v1alpha1.ValueType) error {
	return &ValueTypeError{t: t}
}

func NewValidationError(name string, cause error) error {
	if cause == nil {
		return nil
	}
	return fmt.Errorf("validation error for value %v: %w", name, cause)
}

func NewConfigMapRefError(source v1alpha1.ObjectKeyValueSource, cause error) error {
	return fmt.Errorf("cannot resolve reference to ConfigMap %v.%v: %w", source.Name, source.Namespace, cause)
}

func NewSecretRefError(source v1alpha1.ObjectKeyValueSource, cause error) error {
	return fmt.Errorf("cannot resolve reference to Secret %v.%v: %w", source.Name, source.Namespace, cause)
}

func NewPackageRefError(source v1alpha1.PackageValueSource, cause error) error {
	return fmt.Errorf("cannot resolve reference to value %v in Package %v: %w", source.Value, source.Name, cause)
}

func NewOptionsError(options []string) error {
	return fmt.Errorf("value must be one of: %v", strings.Join(options, ", "))
}

func NewFormatError(format string, cause error) error {
	return fmt.Errorf("value must be a %v: %w", format, cause)
}

func NewPatternError(pattern string) error {
	return fmt.Errorf("value must match '%v'", pattern)
}

var (
	ErrNoDef               = errors.New("no value definition found")
	ErrConstraint          = errors.New("constraint violation")
	ErrConstraintRequired  = fmt.Errorf("%w: Required", ErrConstraint)
	ErrConstraintMin       = fmt.Errorf("%w: Min", ErrConstraint)
	ErrConstraintMax       = fmt.Errorf("%w: Max", ErrConstraint)
	ErrConstraintMinLength = fmt.Errorf("%w: MinLength", ErrConstraint)
	ErrConstraintMaxLength = fmt.Errorf("%w: MaxLength", ErrConstraint)
)
