package clientutils

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
)

func AutoUpdateString(pkg *v1alpha1.Package, disabledStr string) string {
	if pkg != nil {
		if pkg.AutoUpdatesEnabled() {
			return "Enabled"
		}
		return disabledStr
	}
	return ""
}
