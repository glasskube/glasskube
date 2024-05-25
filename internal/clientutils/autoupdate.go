package clientutils

import (
	"strconv"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

func AutoUpdateString(pkg *v1alpha1.Package, disabledStr string) string {
	if pkg != nil {
		if pkg.Annotations != nil {
			autoUpdateValue, ok := pkg.Annotations["packages.glasskube.dev/auto-update"]
			autoUpdateBool, _ := strconv.ParseBool(autoUpdateValue)
			if ok && autoUpdateBool {
				return "Enabled"
			}
		}
		return disabledStr
	}
	return ""
}
