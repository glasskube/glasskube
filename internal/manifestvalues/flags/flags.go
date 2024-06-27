package flags

import (
	"fmt"
	"maps"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/spf13/cobra"
)

type valuesOptionsConfigurer = func(opts *ValuesOptions)

type ValuesOptions struct {
	Values               []string
	KeepOldValues        bool
	KeepOldValuesDefault *bool
}

func NewOptions(conf ...valuesOptionsConfigurer) ValuesOptions {
	var opt ValuesOptions
	for _, fn := range conf {
		fn(&opt)
	}
	return opt
}

var WithKeepOldValuesFlag valuesOptionsConfigurer = func(opts *ValuesOptions) {
	tmp := true
	opts.KeepOldValuesDefault = &tmp
}

func (opts *ValuesOptions) IsValuesSet() bool {
	return (opts.KeepOldValuesDefault != nil && opts.KeepOldValues != *opts.KeepOldValuesDefault) || len(opts.Values) > 0
}

func (opts *ValuesOptions) AddFlagsToCommand(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringArrayVar(&opts.Values, "value", opts.Values,
		"set a value via flag (can be used multiple times).\n"+
			"You can create values referencing data in other resources using the following syntax: "+
			"$<ReferenceKind>$[<specifier>[,<specifier>...]].\n"+
			"For example:\n"+
			" * Reference a ConfigMap key: --value \"name=$ConfigMapRef$namespace,name,key\"\n"+
			" * Reference a Secret key: --value \"name=$SecretRef$namespace,name,key\"\n"+
			" * Reference another Package value: --value \"name=$PackageRef$name,value\"\n")
	if opts.KeepOldValuesDefault != nil {
		flags.BoolVar(&opts.KeepOldValues, "keep-old-values", *opts.KeepOldValuesDefault,
			"set this to false in order to erase any values not specified via --value")
	}
}

func (opts *ValuesOptions) ParseValues(
	oldValues map[string]v1alpha1.ValueConfiguration,
) (map[string]v1alpha1.ValueConfiguration, error) {
	newValues := make(map[string]v1alpha1.ValueConfiguration)
	if opts.KeepOldValues {
		maps.Copy(newValues, oldValues)
	}
	for _, s := range opts.Values {
		split := strings.SplitN(s, "=", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid value format: %v", s)
		}
		key, value := split[0], split[1]

		var valueConfiguration v1alpha1.ValueConfiguration
		if strings.HasPrefix(value, "$ConfigMapRef$") {
			if source, err := parseObjectKeyValueSource(value, "$ConfigMapRef$"); err != nil {
				return nil, fmt.Errorf("value %v is invalid: %v", key, err)
			} else {
				valueConfiguration.ValueFrom = &v1alpha1.ValueReference{ConfigMapRef: source}
			}
		} else if strings.HasPrefix(value, "$SecretRef$") {
			if source, err := parseObjectKeyValueSource(value, "$SecretRef$"); err != nil {
				return nil, fmt.Errorf("value %v is invalid: %v", key, err)
			} else {
				valueConfiguration.ValueFrom = &v1alpha1.ValueReference{SecretRef: source}
			}
		} else if strings.HasPrefix(value, "$PackageRef$") {
			if source, err := parsePackageValueSource(value); err != nil {
				return nil, fmt.Errorf("value %v is invalid: %v", key, err)
			} else {
				valueConfiguration.ValueFrom = &v1alpha1.ValueReference{PackageRef: source}
			}
		} else {
			valueConfiguration.Value = &value
		}
		newValues[key] = valueConfiguration
	}
	return newValues, nil
}

func parseObjectKeyValueSource(value, prefix string) (*v1alpha1.ObjectKeyValueSource, error) {
	if parts, err := parseSourceParts(value, prefix, 3); err != nil {
		return nil, err
	} else {
		return &v1alpha1.ObjectKeyValueSource{Namespace: parts[0], Name: parts[1], Key: parts[2]}, nil
	}
}

func parsePackageValueSource(value string) (*v1alpha1.PackageValueSource, error) {
	if parts, err := parseSourceParts(value, "$PackageRef$", 2); err != nil {
		return nil, err
	} else {
		return &v1alpha1.PackageValueSource{Name: parts[0], Value: parts[1]}, nil
	}
}

func parseSourceParts(value, prefix string, n int) ([]string, error) {
	if parts := strings.SplitN(strings.TrimPrefix(value, prefix), ",", n); len(parts) != n {
		return nil, fmt.Errorf("%v requires %v parameters, got %v", prefix, n, len(parts))
	} else {
		return parts, nil
	}
}
