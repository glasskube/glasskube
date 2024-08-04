package refresh

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

type RefreshTriggerHeaderOnly bool

const (
	RefreshTriggerHeader RefreshTriggerHeaderOnly = true
	RefreshTriggerAll    RefreshTriggerHeaderOnly = false
)

const scopeClusterPackage = "clusterpackage"
const scopePackage = "package"
const segmentHeader = "header"
const RefreshPackageOverview = "refresh-package-overview"
const RefreshClusterPackageOverview = "refresh-clusterpackage-overview"

// GetPackageRefreshDetailId returns the refresh id for the package detail page (or only its header). It is meant
// to be called in situations when there is no manifest at hand, and we therefore only know by the packages' type,
// of which scope it is. This means the pkg has to be installed.
func GetPackageRefreshDetailId(pkg ctrlpkg.Package, headerOnly RefreshTriggerHeaderOnly) string {
	segment := ""
	if headerOnly {
		segment = segmentHeader
	}

	var id string
	var scope string
	if !pkg.IsNamespaceScoped() {
		id = pkg.GetName()
		scope = scopeClusterPackage
	} else if !pkg.IsNil() {
		id = getNamespacedNameId(pkg)
		scope = scopePackage
	}
	return getRefreshId(scope, segment, id)
}

// PackageRefreshDetailId the refresh id for the package detail page (or only its header). It is meant to be called
// in situations where the manifest is at hand (e.g. during template rendering).
func PackageRefreshDetailId(manifest *v1alpha1.PackageManifest, pkg ctrlpkg.Package) string {
	scope, id := getScopeAndId(manifest, pkg)
	return getRefreshId(scope, "", id)
}

// PackageRefreshDetailHeaderId is like PackageRefreshDetailId but only for the header component.
func PackageRefreshDetailHeaderId(manifest *v1alpha1.PackageManifest, pkg ctrlpkg.Package) string {
	scope, id := getScopeAndId(manifest, pkg)
	return getRefreshId(scope, segmentHeader, id)
}

func PackageOverviewRefreshId() string {
	return RefreshPackageOverview
}

func ClusterPackageOverviewRefreshId() string {
	return RefreshClusterPackageOverview
}

func getScopeAndId(manifest *v1alpha1.PackageManifest, pkg ctrlpkg.Package) (string, string) {
	if manifest.Scope.IsCluster() {
		return scopeClusterPackage, manifest.Name
	} else if !pkg.IsNil() {
		return scopePackage, getNamespacedNameId(pkg)
	} else {
		return scopePackage, "---"
	}
}

func getNamespacedNameId(pkg ctrlpkg.Package) string {
	return fmt.Sprintf("%s-%s", pkg.GetNamespace(), pkg.GetName())
}

func getRefreshId(scope, segment, id string) string {
	if segment != "" {
		segment = "-" + segment
	}
	return fmt.Sprintf("refresh-%s-detail%s-%s", scope, segment, id)
}
