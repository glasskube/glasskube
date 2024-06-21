package watch

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

type ownedMapperFunc func(pkg ctrlpkg.Package) []v1alpha1.OwnedResourceRef

var _ ownedMapperFunc = OwnedPackageInfos
var _ ownedMapperFunc = OwnedPackages

func OwnedPackageInfos(pkg ctrlpkg.Package) []v1alpha1.OwnedResourceRef {
	return pkg.GetStatus().OwnedPackageInfos
}

func OwnedPackages(pkg ctrlpkg.Package) []v1alpha1.OwnedResourceRef {
	return pkg.GetStatus().OwnedPackages
}
