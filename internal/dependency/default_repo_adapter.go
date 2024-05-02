package dependency

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	repotypes "github.com/glasskube/glasskube/internal/repo/types"
)

type defaultRepoAdapter struct {
	client repoclient.RepoClientset
}

func (a *defaultRepoAdapter) GetVersions(name string) ([]string, error) {
	packageRepo, err := a.getRepoForPackage(name)
	if err != nil {
		return nil, err
	}
	var idx repotypes.PackageIndex
	if err := a.client.ForRepo(*packageRepo).FetchPackageIndex(name, &idx); err != nil {
		return nil, err
	}
	versions := make([]string, len(idx.Versions))
	for i, item := range idx.Versions {
		versions[i] = item.Version
	}
	return versions, nil
}

func (a *defaultRepoAdapter) GetManifest(name string, version string) (*v1alpha1.PackageManifest, error) {
	if repo, err := a.getRepoForPackage(name); err != nil {
		return nil, err
	} else {
		var manifest v1alpha1.PackageManifest
		return &manifest, a.client.ForRepo(*repo).FetchPackageManifest(name, version, &manifest)
	}
}

func (a *defaultRepoAdapter) getRepoForPackage(name string) (*v1alpha1.PackageRepository, error) {
	if repos, err := a.client.Aggregate().GetReposForPackage(name); err != nil {
		return nil, err
	} else {
		switch len(repos) {
		case 0:
			return nil, fmt.Errorf("%v is not available in any repository", name)
		case 1:
			return &repos[0], nil
		default:
			return nil, fmt.Errorf("%v is available from %v repositories (currently unsupported)", name, len(repos))
		}
	}
}
