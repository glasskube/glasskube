package manifestvalues

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

func ValueAsString(value v1alpha1.ValueConfiguration) string {
	if value.ValueFrom != nil {
		if value.ValueFrom.ConfigMapRef != nil {
			return fmt.Sprintf("reference to '%v' in ConfigMap %v in namespace %v",
				value.ValueFrom.ConfigMapRef.Key, value.ValueFrom.ConfigMapRef.Name, value.ValueFrom.ConfigMapRef.Namespace)
		} else if value.ValueFrom.SecretRef != nil {
			return fmt.Sprintf("reference to '%v' in Secret %v in namespace %v",
				value.ValueFrom.SecretRef.Key, value.ValueFrom.SecretRef.Name, value.ValueFrom.SecretRef.Namespace)
		} else if value.ValueFrom.PackageRef != nil {
			return fmt.Sprintf("reference to value '%v' of Package %v",
				value.ValueFrom.PackageRef.Value, value.ValueFrom.PackageRef.Name)
		}
	} else if value.Value != nil {
		return *value.Value
	}
	return "n/a"
}
