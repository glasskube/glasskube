package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

const (
	namespaceKey = "namespace"
	nameKey      = "name"
	keyKey       = "key"
	packageKey   = "package"
	valueKey     = "value"
	refKindKey   = "refKind"
)

func extractValues(r *http.Request, manifest *v1alpha1.PackageManifest) (map[string]v1alpha1.ValueConfiguration, error) {
	values := make(map[string]v1alpha1.ValueConfiguration)
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	for valueName, valueDef := range manifest.ValueDefinitions {
		if refKindVal := r.Form.Get(fmt.Sprintf("%s[%s]", valueName, refKindKey)); refKindVal == "ConfigMap" {
			values[valueName] = v1alpha1.ValueConfiguration{
				ValueFrom: &v1alpha1.ValueReference{
					ConfigMapRef: extractObjectKeyValueSource(r, valueName),
				},
			}
		} else if refKindVal == "Secret" {
			values[valueName] = v1alpha1.ValueConfiguration{
				ValueFrom: &v1alpha1.ValueReference{
					SecretRef: extractObjectKeyValueSource(r, valueName),
				},
			}
		} else if refKindVal == "Package" {
			values[valueName] = v1alpha1.ValueConfiguration{
				ValueFrom: &v1alpha1.ValueReference{
					PackageRef: extractPackageValueSource(r, valueName),
				},
			}
		} else if refKindVal == "" {
			formVal := r.Form.Get(valueName)
			if valueDef.Type == v1alpha1.ValueTypeBoolean {
				boolStr := strconv.FormatBool(false)
				if strings.ToLower(formVal) == "on" {
					boolStr = strconv.FormatBool(true)
				}
				values[valueName] = v1alpha1.ValueConfiguration{Value: &boolStr}
			} else {
				values[valueName] = v1alpha1.ValueConfiguration{Value: &formVal}
			}
		} else {
			return nil, fmt.Errorf("cannot extract value %v because of unknown reference kind %v", valueName, refKindVal)
		}
	}
	return values, nil
}

func extractObjectKeyValueSource(r *http.Request, valueName string) *v1alpha1.ObjectKeyValueSource {
	namespaceFormKey := fmt.Sprintf("%s[%s]", valueName, namespaceKey)
	nameFormKey := fmt.Sprintf("%s[%s]", valueName, nameKey)
	keyFormKey := fmt.Sprintf("%s[%s]", valueName, keyKey)
	return &v1alpha1.ObjectKeyValueSource{
		Name:      r.Form.Get(nameFormKey),
		Namespace: r.Form.Get(namespaceFormKey),
		Key:       r.Form.Get(keyFormKey),
	}
}

func extractPackageValueSource(r *http.Request, valueName string) *v1alpha1.PackageValueSource {
	packageFormKey := fmt.Sprintf("%s[%s]", valueName, packageKey)
	valueFormKey := fmt.Sprintf("%s[%s]", valueName, valueKey)
	return &v1alpha1.PackageValueSource{
		Name:  r.Form.Get(packageFormKey),
		Value: r.Form.Get(valueFormKey),
	}
}
