package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo/types"
)

type PackageRepoIndexFetcher interface {
	FetchPackageRepoIndex(target *types.PackageRepoIndex) error
}

type LatestVersionGetter interface {
	GetLatestVersion(pkgName string) (string, error)
}

type RepoClient interface {
	PackageRepoIndexFetcher
	LatestVersionGetter
	FetchLatestPackageManifest(name string, target *packagesv1alpha1.PackageManifest) (version string, err error)
	FetchPackageManifest(name, version string, target *packagesv1alpha1.PackageManifest) error
	FetchPackageIndex(name string, target *types.PackageIndex) error
	GetPackageManifestURL(name, version string) (string, error)
}

type RepoAggregator interface {
	PackageRepoIndexFetcher
	LatestVersionGetter
	GetReposForPackage(name string) ([]packagesv1alpha1.PackageRepository, error)
}

type RepoClientset interface {
	ForPackage(pkg v1alpha1.Package) RepoClient
	ForRepoWithName(name string) RepoClient
	ForRepo(repo v1alpha1.PackageRepository) RepoClient
	Default() RepoClient
	Aggregate() RepoAggregator
}
