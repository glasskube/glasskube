package clientutils

import (
	"reflect"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

func AutoUpdateString(pkg ctrlpkg.Package, disabledStr string) string {
	if !reflect.ValueOf(pkg).IsNil() {
		if pkg.AutoUpdatesEnabled() {
			return "Enabled"
		}
		return disabledStr
	}
	return ""
}
