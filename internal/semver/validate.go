package semver

import "github.com/Masterminds/semver/v3"

func ValidateConstraint(version, constraint string) error {
	if parsedVersion, err := semver.NewVersion(version); err != nil {
		return err
	} else if parsedConstraint, err := semver.NewConstraint(constraint); err != nil {
		return err
	} else {
		return ValidateVersionConstraint(parsedVersion, parsedConstraint)
	}
}

func ValidateVersionConstraint(version *semver.Version, constraint *semver.Constraints) error {
	_, err := constraint.Validate(version)
	return ErrConstraintValidation(err)
}
