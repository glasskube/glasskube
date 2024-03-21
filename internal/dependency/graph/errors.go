package graph

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func ErrNotInstalled(name string) error {
	err := NotInstalledError(name)
	return &err
}

// NotInstalledError indicates that a package is missing
type NotInstalledError string

func (err *NotInstalledError) Error() string {
	return fmt.Sprintf("%v not installed", string(*err))
}

func ErrConstraint(name string, version *semver.Version, constraint *semver.Constraints, cause error) error {
	return &ConstraintError{Name: name, Version: version, Constraint: constraint, cause: cause}
}

// ConstraintError indicates that a constraint has been violated
type ConstraintError struct {
	Name       string
	Version    *semver.Version
	Constraint *semver.Constraints
	cause      error
}

func (err *ConstraintError) Error() string {
	return fmt.Sprintf("constraint %v violated: %v", err.Constraint, err.cause)
}

func ErrDependency(name, dep string, cause error) error {
	return &DependencyError{Name: name, Dependency: dep, cause: cause}
}

// DependencyError idicates that a dependency is not met
type DependencyError struct {
	Name, Dependency string
	cause            error
}

func (err *DependencyError) Error() string {
	return fmt.Sprintf("unmet dependency %v -> %v: %v", err.Name, err.Dependency, err.cause)
}

func (err *DependencyError) Is(other error) bool {
	_, ok := other.(*DependencyError)
	return ok
}

func (err *DependencyError) Unwrap() error {
	return err.cause
}

var _ error = &ConstraintError{}
var _ error = &DependencyError{}
