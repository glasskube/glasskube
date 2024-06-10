package names

import (
	"fmt"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

func PackageInfoName(pkg ctrlpkg.PackageCommon) string {
	spec := pkg.GetSpec()
	return escapeResourceName(fmt.Sprintf("%v--%v", spec.PackageInfo.Name, spec.PackageInfo.Version))
}
