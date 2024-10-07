package dependency

import (
	"fmt"

	repoerror "github.com/glasskube/glasskube/internal/repo/error"

	"go.uber.org/multierr"

	"github.com/glasskube/glasskube/api/v1alpha1"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	repotypes "github.com/glasskube/glasskube/internal/repo/types"
)

type defaultRepoAdapter struct {
	client repoclient.RepoClientset
}

func (a *defaultRepoAdapter) GetVersions(name string) ([]string, error) {
	packageRepo, repoErr := a.getRepoForPackage(name)
	if repoerror.IsComplete(repoErr) {
		return nil, repoErr
	}
	var idx repotypes.PackageIndex
	if err := a.client.ForRepo(*packageRepo).FetchPackageIndex(name, &idx); err != nil {
		return nil, multierr.Append(repoErr, err)
	}
	versions := make([]string, len(idx.Versions))
	for i, item := range idx.Versions {
		versions[i] = item.Version
	}
	return versions, repoErr
}

func (a *defaultRepoAdapter) GetManifest(name string, version string) (*v1alpha1.PackageManifest, error) {
	if repo, err := a.getRepoForPackage(name); repoerror.IsComplete(err) {
		return nil, err
	} else {
		var manifest v1alpha1.PackageManifest
		return &manifest, multierr.Append(a.client.ForRepo(*repo).FetchPackageManifest(name, version, &manifest), err)
	}
}

func (a *defaultRepoAdapter) getRepoForPackage(name string) (*v1alpha1.PackageRepository, error) {
	repos, err := a.client.Meta().GetReposForPackage(name)
	switch len(repos) {
	case 0:
		return nil, multierr.Append(fmt.Errorf("\"%v\" is not available in any repository", name), err)
	case 1:
		return &repos[0], err
	default:
		return nil, multierr.Append(err,
			fmt.Errorf("%v is available from %v repositories (currently unsupported)", name, len(repos)),
		)
	}
}

func (a *defaultRepoAdapter) GetManifestFromRepo(name string, version string, repositoryName string) (*v1alpha1.PackageManifest, error) {
	var manifest v1alpha1.PackageManifest
	if err := a.client.ForRepoWithName(repositoryName).FetchPackageManifest(name, version, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}
