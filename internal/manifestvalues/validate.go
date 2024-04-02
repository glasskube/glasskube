package manifestvalues

import (
	"regexp"
	"strconv"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"go.uber.org/multierr"
)

type validateFn func(def v1alpha1.ValueDefinition, value string) error

var (
	validateMinLength validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.MinLength > 0 && len(value) < def.Constraints.MinLength {
			return ErrConstraintMinLength
		}
		return nil
	}

	validateMaxLength validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.MaxLength > 0 && len(value) > def.Constraints.MaxLength {
			return ErrConstraintMaxLength
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
		if def.Constraints.Min > 0 {
			if i, err := strconv.Atoi(value); err != nil {
				return err
			} else if i < def.Constraints.Min {
				return ErrConstraintMin
			}
		}
		return nil
	}

	validateMax validateFn = func(def v1alpha1.ValueDefinition, value string) error {
		if def.Constraints.Max > 0 {
			if i, err := strconv.Atoi(value); err != nil {
				return err
			} else if i > def.Constraints.Max {
				return ErrConstraintMax
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
		if len(def.Constraints.Pattern) > 0 {
			if re, err := regexp.Compile(def.Constraints.Pattern); err != nil {
				return err
			} else if !re.MatchString(value) {
				return NewPatternError(def.Constraints.Pattern)
			}
		}
		return nil
	}
)

type validateFns []validateFn

func (v *validateFns) Validate(def v1alpha1.ValueDefinition, value string) (err error) {
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

func ValidateResolvedValues(manifest v1alpha1.PackageManifest, values map[string]string) (err error) {
	namesFromDef := make(map[string]struct{})
	for name, def := range manifest.ValueDefinitions {
		namesFromDef[name] = struct{}{}
		if validators, ok := validatorsForType[def.Type]; !ok {
			multierr.AppendInto(&err, NewValidationError(name, NewValueTypeError(def.Type)))
		} else if value, ok := values[name]; ok {
			multierr.AppendInto(&err, NewValidationError(name, validators.Validate(def, value)))
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
