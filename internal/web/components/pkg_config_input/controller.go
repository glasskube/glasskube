package pkg_config_input

import (
	"fmt"
	"strconv"

	"github.com/glasskube/glasskube/internal/web/util"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"

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
	PackageHref        string
}

func getStringValue(pkg ctrlpkg.Package, valueName string, valueDefinition *v1alpha1.ValueDefinition) string {
	if !pkg.IsNil() {
		if valueConfiguration, ok := pkg.GetSpec().Values[valueName]; ok {
			if valueConfiguration.Value != nil {
				return *valueConfiguration.Value
			}
		}
	}
	return valueDefinition.DefaultValue
}

func getBoolValue(pkg ctrlpkg.Package, valueName string, valueDefinition *v1alpha1.ValueDefinition) bool {
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

func getExistingReferenceAndKind(pkg ctrlpkg.Package, valueName string) (*v1alpha1.ValueReference, string) {
	if !pkg.IsNil() {
		if val, ok := pkg.GetSpec().Values[valueName]; ok {
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

func getOrCreateReference(
	pkg ctrlpkg.Package, valueName string, desiredRefKind *string) (v1alpha1.ValueReference, string) {
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
	pkg ctrlpkg.Package,
	repositoryName string,
	selectedVersion string,
	manifest *v1alpha1.PackageManifest,
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
		PkgName:            manifest.Name,
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
		PackageHref:        util.GetPackageHrefWithFallback(pkg, manifest),
	}
}
