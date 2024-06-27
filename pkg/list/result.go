package list

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo/types"
)

type result struct {
	IndexItem      *types.MetaIndexItem
	ClusterPackage *v1alpha1.ClusterPackage
	Packages       []*v1alpha1.Package
	PackageInfo    *v1alpha1.PackageInfo
	Repositories   []*v1alpha1.PackageRepository
}

func (item result) Installed() bool {
	return item.ClusterPackage != nil
}

func (item result) Outdated() bool {
	// TODO check again whether this works correctly in relation with multiple repos?
	return item.ClusterPackage != nil && item.IndexItem != nil &&
		item.ClusterPackage.Spec.PackageInfo.Version != item.IndexItem.LatestVersion
}
