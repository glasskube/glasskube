package graph

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func ErrNotInstalled(ref PackageRef) error {
	err := NotInstalledError(ref)
	return &err
}

// NotInstalledError indicates that a package is missing
type NotInstalledError PackageRef

func (err *NotInstalledError) Error() string {
	return fmt.Sprintf("%v not installed", err.PackageName)
}

func ErrConstraint(pkgRef PackageRef, version *semver.Version, constraint *semver.Constraints, cause error) error {
	return &ConstraintError{Package: pkgRef, Version: version, Constraint: constraint, cause: cause}
}

// ConstraintError indicates that a constraint has been violated
type ConstraintError struct {
	Package    PackageRef
	Version    *semver.Version
	Constraint *semver.Constraints
	cause      error
}

func (err *ConstraintError) Error() string {
	return fmt.Sprintf("constraint %v violated: %v", err.Constraint, err.cause)
}

func ErrDependency(ref, dep PackageRef, cause error) error {
	return &DependencyError{Package: ref, Dependency: dep, cause: cause}
}

// DependencyError idicates that a dependency is not met
type DependencyError struct {
	Package, Dependency PackageRef
	cause               error
}

func (err *DependencyError) Error() string {
	return fmt.Sprintf("unmet dependency %v -> %v: %v", err.Package.PackageName, err.Dependency.PackageName, err.cause)
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
