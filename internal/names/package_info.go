package names

import (
	"strings"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

func PackageInfoName(pkg ctrlpkg.Package) string {
	spec := pkg.GetSpec()
	parts := []string{spec.PackageInfo.Name, spec.PackageInfo.Version}
	if spec.PackageInfo.RepositoryName != "" {
		parts = append([]string{spec.PackageInfo.RepositoryName}, parts...)
	}
	return escapeResourceName(strings.Join(parts, "--"))
}
