package pkg_config_input

import (
	"fmt"
	"strconv"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

type PkgConfigInputRenderOptions struct {
	Autofocus      bool
	DesiredRefKind *string
}

type PkgConfigInputDatalistOptions struct {
	Namespaces []string
	Names      []string
	Keys       []string
}

type pkgConfigInputInput struct {
	RepositoryName     string
	SelectedVersion    string
	PkgName            string
	ValueName          string
	ValueDefinition    v1alpha1.ValueDefinition
	StringValue        string
	BoolValue          bool
	FormLabel          string
	FormId             string
	ContainerId        string
	ValueReference     v1alpha1.ValueReference
	ValueReferenceKind string
	ValueError         error
	Autofocus          bool
	DatalistOptions    *PkgConfigInputDatalistOptions
}

func getStringValue(pkg *v1alpha1.Package, valueName string, valueDefinition *v1alpha1.ValueDefinition) string {
	if pkg != nil {
		if valueConfiguration, ok := pkg.Spec.Values[valueName]; ok {
			if valueConfiguration.Value != nil {
				return *valueConfiguration.Value
			}
		}
	}
	return valueDefinition.DefaultValue
}

func getBoolValue(pkg *v1alpha1.Package, valueName string, valueDefinition *v1alpha1.ValueDefinition) bool {
	if valueDefinition.Type == v1alpha1.ValueTypeBoolean {
		strVal := getStringValue(pkg, valueName, valueDefinition)
		if valBool, err := strconv.ParseBool(strVal); err == nil {
			return valBool
		}
	}
	return false
}

func getLabel(valueName string, valueDefinition *v1alpha1.ValueDefinition) string {
	inputLabel := valueName
	if valueDefinition.Metadata.Label != "" {
		inputLabel = valueDefinition.Metadata.Label
	}
	return inputLabel
}

func getExistingReferenceAndKind(pkg *v1alpha1.Package, valueName string) (*v1alpha1.ValueReference, string) {
	if pkg != nil {
		if val, ok := pkg.Spec.Values[valueName]; ok {
			if val.Value == nil && val.ValueFrom != nil {
				if val.ValueFrom.ConfigMapRef != nil {
					return val.ValueFrom, "ConfigMap"
				} else if val.ValueFrom.SecretRef != nil {
					return val.ValueFrom, "Secret"
				} else if val.ValueFrom.PackageRef != nil {
					return val.ValueFrom, "Package"
				}
			}
		}
	}
	return nil, ""
}

func getOrCreateReference(pkg *v1alpha1.Package, valueName string, desiredRefKind *string) (v1alpha1.ValueReference, string) {
	existingReference, existingRefKind := getExistingReferenceAndKind(pkg, valueName)
	if desiredRefKind != nil && *desiredRefKind != existingRefKind {
		return v1alpha1.ValueReference{}, *desiredRefKind
	} else if existingReference != nil {
		return *existingReference, existingRefKind
	} else {
		return v1alpha1.ValueReference{}, existingRefKind
	}
}

func ForPkgConfigInput(
	pkg *v1alpha1.Package,
	repositoryName string,
	selectedVersion string,
	pkgName string,
	valueName string,
	valueDefinition v1alpha1.ValueDefinition,
	valueError error,
	datalistOptions *PkgConfigInputDatalistOptions,
	options *PkgConfigInputRenderOptions,
) *pkgConfigInputInput {
	if options == nil {
		options = &PkgConfigInputRenderOptions{}
	}
	valueReference, valueReferenceKind := getOrCreateReference(pkg, valueName, options.DesiredRefKind)
	return &pkgConfigInputInput{
		RepositoryName:     repositoryName,
		SelectedVersion:    selectedVersion,
		PkgName:            pkgName,
		ValueName:          valueName,
		ValueDefinition:    valueDefinition,
		StringValue:        getStringValue(pkg, valueName, &valueDefinition),
		BoolValue:          getBoolValue(pkg, valueName, &valueDefinition),
		FormLabel:          getLabel(valueName, &valueDefinition),
		FormId:             fmt.Sprintf("input-%v", valueName),
		ContainerId:        fmt.Sprintf("input-container-%v", valueName),
		ValueReference:     valueReference,
		ValueReferenceKind: valueReferenceKind,
		ValueError:         valueError,
		Autofocus:          options.Autofocus,
		DatalistOptions:    datalistOptions,
	}
}
