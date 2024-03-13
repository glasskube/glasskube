package client

import (
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo/types"
)

type RepoClient interface {
	FetchLatestPackageManifest(repoURL, name string, target *packagesv1alpha1.PackageManifest) (version string, err error)
	FetchPackageManifest(repoURL, name, version string, target *packagesv1alpha1.PackageManifest) error
	FetchPackageIndex(repoURL, name string, target *types.PackageIndex) error
	FetchPackageRepoIndex(repoURL string, target *types.PackageRepoIndex) error
	GetLatestVersion(repoURL string, pkgName string) (string, error)
	GetPackageManifestURL(repoURL, name, version string) (string, error)
}
