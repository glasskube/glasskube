package util

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

func GetPackageHref(pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest) string {
	var pkgHref string
	if manifest.Scope == nil || *manifest.Scope == v1alpha1.ScopeCluster {
		// Scope == nil is the fallback for all older packages – it will only be wrong for quickwit (the first non-cluster
		// package), and only when someone selects an outdated version
		pkgHref = fmt.Sprintf("/clusterpackages/%s", manifest.Name)
	} else {
		pkgPath := ""
		if !pkg.IsNil() {
			pkgPath = fmt.Sprintf("/%s/%s", pkg.GetNamespace(), pkg.GetName())
		} else {
			// TODO is this correct in every case?
			pkgPath = "/-/-" // not installed yet
		}
		pkgHref = fmt.Sprintf("/packages/%s%s", manifest.Name, pkgPath)
	}
	return pkgHref
}
