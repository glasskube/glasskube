package util

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

func GetPackageHref(pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest) string {
	return getPackageHref(pkg, manifest, false)
}

func GetPackageHrefWithFallback(pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest) string {
	return getPackageHref(pkg, manifest, true)
}

func getPackageHref(pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest, withFallback bool) string {
	if manifest.Scope.IsCluster() {
		return GetClusterPkgHref(manifest.Name)
	} else {
		if !pkg.IsNil() {
			return GetNamespacedPkgHref(manifest.Name, pkg.GetNamespace(), pkg.GetName())
		} else if withFallback {
			return GetNamespacedPkgHref(manifest.Name, "-", "-") // not installed yet
		}
		return GetNamespacedPkgHref(manifest.Name, "", "")
	}
}

func GetClusterPkgHref(pkgName string) string {
	return fmt.Sprintf("/clusterpackages/%s", pkgName)
}

func GetNamespacedPkgHref(manifestName string, namespace string, name string) string {
	if namespace != "" && name != "" {
		return fmt.Sprintf("/packages/%s/%s/%s", manifestName, namespace, name)
	} else {
		return fmt.Sprintf("/packages/%s", manifestName)
	}
}
