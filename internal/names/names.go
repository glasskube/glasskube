package names

import (
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
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

func HelmResourceName(pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest) string {
	return strings.Join([]string{pkg.GetName(), manifest.Name}, "-")
}
