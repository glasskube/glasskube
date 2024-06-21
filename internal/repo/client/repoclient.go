package client

import (
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/repo/types"
)

type LatestVersionGetter interface {
	GetLatestVersion(pkgName string) (string, error)
}

type RepoClient interface {
	LatestVersionGetter
	FetchPackageRepoIndex(target *types.PackageRepoIndex) error
	FetchLatestPackageManifest(name string, target *packagesv1alpha1.PackageManifest) (version string, err error)
	FetchPackageManifest(name, version string, target *packagesv1alpha1.PackageManifest) error
	FetchPackageIndex(name string, target *types.PackageIndex) error
	GetPackageManifestURL(name, version string) (string, error)
}

type RepoMetaclient interface {
	LatestVersionGetter
	FetchMetaIndex(target *types.MetaIndex) error
	GetReposForPackage(name string) ([]packagesv1alpha1.PackageRepository, error)
}

type RepoClientset interface {
	ForPackage(pkg ctrlpkg.Package) RepoClient
	ForRepoWithName(name string) RepoClient
	ForRepo(repo packagesv1alpha1.PackageRepository) RepoClient
	Default() RepoClient
	Meta() RepoMetaclient
}
