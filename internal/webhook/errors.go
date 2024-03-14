package webhook

import (
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/dependency"
)

var ErrInvalidObject = errors.New("validator called with unexpected object type")
var ErrDependencyConflict = errors.New("dependency conflict")
var ErrTransitiveDependency = errors.New("support for transitive dependencies is not implemented yet")

func newConflictError(conflicts dependency.Conflicts) error {
	return fmt.Errorf("%w: %v", ErrDependencyConflict, conflicts)
}

func newConflictErrorDelete(owner string) error {
	return fmt.Errorf("%w: dependency of %v", ErrDependencyConflict, owner)
}

func newTransitiveError(requirement dependency.PackageWithVersion, dependency v1alpha1.Dependency) error {
	return fmt.Errorf("%w: requirement %v (%v) depends on %v (%v)",
		ErrTransitiveDependency, requirement.Name, requirement.Version, dependency.Name, dependency.Version)
}

func isErrDependencyConflict(err error) bool   { return errors.Is(err, ErrDependencyConflict) }
func isErrTransitiveDependency(err error) bool { return errors.Is(err, ErrTransitiveDependency) }
