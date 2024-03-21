package dependency

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/repo/client"
)

type defaultRepoAdapter struct {
	repo client.RepoClient
}

func (a *defaultRepoAdapter) GetVersions(repoURL string, name string) ([]string, error) {
	var idx repo.PackageIndex
	if err := a.repo.FetchPackageIndex(repoURL, name, &idx); err != nil {
		return nil, err
	}
	versions := make([]string, len(idx.Versions))
	for i, item := range idx.Versions {
		versions[i] = item.Version
	}
	return versions, nil
}

func (a *defaultRepoAdapter) GetManifest(repoURL string, name string, version string) (*v1alpha1.PackageManifest, error) {
	var manifest v1alpha1.PackageManifest
	return &manifest, a.repo.FetchPackageManifest(repoURL, name, version, &manifest)
}
