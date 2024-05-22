package list

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo/types"
)

type result struct {
	IndexItem    *types.MetaIndexItem
	Package      *v1alpha1.Package
	PackageInfo  *v1alpha1.PackageInfo
	Repositories []*v1alpha1.PackageRepository
}

func (item result) Installed() bool {
	return item.Package != nil
}

func (item result) Outdated() bool {
	return item.Package != nil && item.IndexItem != nil &&
		item.Package.Spec.PackageInfo.Version != item.IndexItem.LatestVersion
}
