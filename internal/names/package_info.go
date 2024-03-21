package names

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

func PackageInfoName(pkg v1alpha1.Package) string {
	return escapeResourceName(fmt.Sprintf("%v--%v", pkg.Spec.PackageInfo.Name, pkg.Spec.PackageInfo.Version))
}
