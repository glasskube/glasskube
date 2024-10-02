package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo/client/auth"
	"github.com/glasskube/glasskube/internal/repo/types"
)

// errorclient is a no-op implementation of RepoClient that defers returning an error to the method call.
// This is done to improve ergonomics of the RepoClientset, such that the ForRepo function does not return an error.
type errorclient struct {
	auth.NoopAuthenticator
	err error
}

var _ RepoClient = &errorclient{}

// FetchLatestPackageManifest implements RepoClient.
func (e *errorclient) FetchLatestPackageManifest(name string, target *v1alpha1.PackageManifest) (version string, err error) {
	return "", e.err
}

// FetchPackageIndex implements RepoClient.
func (e *errorclient) FetchPackageIndex(name string, target *types.PackageIndex) error {
	return e.err
}

// FetchPackageManifest implements RepoClient.
func (e *errorclient) FetchPackageManifest(name string, version string, target *v1alpha1.PackageManifest) error {
	return e.err
}

// FetchPackageRepoIndex implements RepoClient.
func (e *errorclient) FetchPackageRepoIndex(target *types.PackageRepoIndex) error {
	return e.err
}

// GetLatestVersion implements RepoClient.
func (e *errorclient) GetLatestVersion(pkgName string) (string, error) {
	return "", e.err
}

// GetPackageManifestURL implements RepoClient.
func (e *errorclient) GetPackageManifestURL(name string, version string) (string, error) {
	return "", e.err
}
