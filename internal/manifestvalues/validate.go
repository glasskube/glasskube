package manifestvalues

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"go.uber.org/multierr"
)

type validateFn func(def v1alpha1.ValueDefinition, value string) error

var (
	validateMinLength validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.MinLength != nil && len(value) < *def.Constraints.MinLength {
			return fmt.Errorf("%w: %v", ErrConstraintMinLength, *def.Constraints.MinLength)
		}
		return nil
	}

	validateMaxLength validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.MaxLength != nil && len(value) > *def.Constraints.MaxLength {
			return fmt.Errorf("%w: %v", ErrConstraintMaxLength, *def.Constraints.MaxLength)
		}
		return nil
	}

	validateFormatNumber validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if _, err := strconv.Atoi(value); err != nil {
			return NewFormatError("number", err)
		}
		return nil
	}

	validateFormatBoolean validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if _, err := strconv.ParseBool(value); err != nil {
			return NewFormatError("boolean", err)
		}
		return nil
	}

	validateMin validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.Min != nil {
			if i, err := strconv.Atoi(value); err != nil {
				return err
			} else if i < *def.Constraints.Min {
				return fmt.Errorf("%w: %v", ErrConstraintMin, *def.Constraints.Min)
			}
		}
		return nil
	}

	validateMax validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.Max != nil {
			if i, err := strconv.Atoi(value); err != nil {
				return err
			} else if i > *def.Constraints.Max {
				return fmt.Errorf("%w: %v", ErrConstraintMax, *def.Constraints.Max)
			}
		}
		return nil
	}

	validateOptions validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		for _, o := range def.Options {
			if value == o {
				return nil
			}
		}
		return NewOptionsError(def.Options)
	}

	validatePattern validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.Pattern != nil {
			if re, err := regexp.Compile(*def.Constraints.Pattern); err != nil {
				return err
			} else if !re.MatchString(value) {
				return NewPatternError(*def.Constraints.Pattern)
			}
		}
		return nil
	}
)

type validateFns []validateFn

func (v *validateFns) validate(def v1alpha1.ValueDefinition, value string) (err error) {
	for _, fn := range *v {
		multierr.AppendInto(&err, fn(def, value))
	}
	return
}

var validatorsForType = map[v1alpha1.ValueType]validateFns{
	v1alpha1.ValueTypeText:    {validateMaxLength, validateMinLength, validatePattern},
	v1alpha1.ValueTypeNumber:  {validateFormatNumber, validateMin, validateMax, validatePattern},
	v1alpha1.ValueTypeOptions: {validateOptions},
	v1alpha1.ValueTypeBoolean: {validateFormatBoolean},
}

func validate(manifest v1alpha1.PackageManifest, values map[string]validationTarget) (err error) {
	namesFromDef := make(map[string]struct{})
	for name, def := range manifest.ValueDefinitions {
		namesFromDef[name] = struct{}{}
		if value, ok := values[name]; ok {
			if !value.Skip() {
				multierr.AppendInto(&err, ValidateSingle(name, def, value.Get()))
			}
		} else if def.Constraints.Required {
			multierr.AppendInto(&err, NewValidationError(name, ErrConstraintRequired))
		}
	}
	for name := range values {
		// check for value configurations that don't have a matching value definition
		if _, ok := namesFromDef[name]; !ok {
			multierr.AppendInto(&err, NewValidationError(name, ErrNoDef))
		}
	}
	return
}

func ValidateSingle(name string, def v1alpha1.ValueDefinition, value string) (err error) {
	if validators, ok := validatorsForType[def.Type]; !ok {
		multierr.AppendInto(&err, NewValidationError(name, NewValueTypeError(def.Type)))
	} else {
		multierr.AppendInto(&err, NewValidationError(name, validators.validate(def, value)))
	}
	return
}

type validationTarget interface {
	Get() string
	Skip() bool
}

type acutalValue string

func (v acutalValue) Get() string {
	return string(v)
}

func (v acutalValue) Skip() bool {
	return false
}

type noopValue struct{}

func (noopValue) Get() string {
	return ""
}

func (noopValue) Skip() bool {
	return true
}

func targetsForResolvedValues(values map[string]string) map[string]validationTarget {
	result := make(map[string]validationTarget)
	for name, value := range values {
		result[name] = acutalValue(value)
	}
	return result
}

func targetsForPackage(pkg ctrlpkg.Package) map[string]validationTarget {
	result := make(map[string]validationTarget)
	for name, value := range pkg.GetSpec().Values {
		if value.Value != nil {
			result[name] = acutalValue(*value.Value)
		} else {
			// reference values should be skipped
			result[name] = noopValue{}
		}
	}
	return result
}

func ValidateResolvedValues(manifest v1alpha1.PackageManifest, values map[string]string) error {
	return validate(manifest, targetsForResolvedValues(values))
}

// ValidatePackage performs a partial validation of a packages value configurations.
// Reference values are **not resolved**, so constraint validation is skipped for
// these values. If instead you want to validate **all** values, please resolve all
// references first, using a Resolver instance, and then use ValidateResolvedValues.
func ValidatePackage(manifest v1alpha1.PackageManifest, pkg ctrlpkg.Package) error {
	return validate(manifest, targetsForPackage(pkg))
}
