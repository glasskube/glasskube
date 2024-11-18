package webhook

import (
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/dependency"
)

var ErrInvalidObject = errors.New("validator called with unexpected object type")
var ErrDependencyConflict = errors.New("dependency conflict")
var ErrPackagesInstalled = errors.New("dependent package(s) installed")

func newConflictError(conflicts dependency.Conflicts) error {
	return fmt.Errorf("%w: %v", ErrDependencyConflict, conflicts)
}

func newConflictErrorDelete(err error) error {
	return fmt.Errorf("%w: %w", ErrDependencyConflict, err)
}

func newErrPackagesInstalled(packageInfoSpecs []v1alpha1.PackageInfoSpec) error {
	return fmt.Errorf("%w: %v", ErrPackagesInstalled, packageInfoSpecs)
}

func isErrDependencyConflict(err error) bool { return errors.Is(err, ErrDependencyConflict) }
