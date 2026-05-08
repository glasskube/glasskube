package semver

import (
	"fmt"
	"regexp"

	"github.com/Masterminds/semver/v3"
)

var constraintBuildMetadataRegex = regexp.MustCompile(`v?[0-9]+(?:\.[0-9]+){0,2}(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*`)

func NewConstraint(constraint string) (*semver.Constraints, error) {
	if constraintBuildMetadataRegex.MatchString(constraint) {
		return nil, fmt.Errorf("build metadata is not supported in version constraints: %s", constraint)
	}
	return semver.NewConstraint(constraint)
}

func ValidateConstraint(version, constraint string) error {
	if parsedVersion, err := semver.NewVersion(version); err != nil {
		return err
	} else if parsedConstraint, err := NewConstraint(constraint); err != nil {
		return err
	} else {
		return ValidateVersionConstraint(parsedVersion, parsedConstraint)
	}
}

func ValidateVersionConstraint(version *semver.Version, constraint *semver.Constraints) error {
	_, err := constraint.Validate(version)
	return ErrConstraintValidation(err)
}
