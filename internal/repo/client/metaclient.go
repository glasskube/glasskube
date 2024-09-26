package client

import (
	"context"
	"slices"

	repoerror "github.com/glasskube/glasskube/internal/repo/error"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/maputils"
	"github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/internal/util"
	"go.uber.org/multierr"
)

type metaclient struct {
	clientset *defaultClientset
}

// FetchPackageRepoIndex implements RepoMetaclient.
func (d metaclient) FetchMetaIndex(target *types.MetaIndex) error {
	if repoList, err := d.clientset.client.ListPackageRepositories(context.TODO()); err != nil {
		return err
	} else {
		var compositeErr error
		indexMap := make(map[string]types.MetaIndexItem)
		util.SortBy(repoList.Items, func(repo v1alpha1.PackageRepository) string { return repo.Name })
		// Reverse the items, because metaItem.PackageRepoIndexItem should be set to the relevant entry from the FIRST
		// repository that is not the default.
		slices.Reverse(repoList.Items)
		for _, repo := range repoList.Items {
			var index types.PackageRepoIndex
			if err := d.clientset.ForRepo(repo).FetchPackageRepoIndex(&index); err != nil {
				multierr.AppendInto(&compositeErr, err)
			} else {
				for _, item := range index.Packages {
					if metaItem, ok := indexMap[item.Name]; !ok {
						indexMap[item.Name] = types.MetaIndexItem{
							PackageRepoIndexItem: item,
							Repos:                []string{repo.Name},
						}
					} else {
						// Insert current repo at the head, because we loop over repos in reverse order
						metaItem.Repos = append([]string{repo.Name}, metaItem.Repos...)
						// The LatestVersion is set to the hightest semver across all repos
						actualLatestVersion := metaItem.LatestVersion
						if semver.IsUpgradable(actualLatestVersion, item.LatestVersion) {
							actualLatestVersion = item.LatestVersion
						}
						if !repo.IsDefaultRepository() {
							metaItem.PackageRepoIndexItem = item
						}
						metaItem.LatestVersion = actualLatestVersion
						indexMap[item.Name] = metaItem
					}
				}
			}
		}
		*target = types.MetaIndex{
			Packages: make([]types.MetaIndexItem, len(indexMap)),
		}
		for i, name := range maputils.KeysSorted(indexMap) {
			target.Packages[i] = indexMap[name]
		}
		return compositeErr
	}
}

// GetLatestVersion implements RepoMetaclient.
func (d metaclient) GetLatestVersion(pkgName string) (string, error) {
	if repoList, err := d.clientset.client.ListPackageRepositories(context.TODO()); err != nil {
		return "", err
	} else {
		var latest string
		for _, repo := range repoList.Items {
			var index types.PackageIndex
			if err := d.clientset.ForRepo(repo).FetchPackageIndex(pkgName, &index); err != nil {
				return "", err
			}
			if latest == "" || semver.IsUpgradable(latest, index.LatestVersion) {
				latest = index.LatestVersion
			}
		}
		return latest, nil
	}
}

// GetReposForPackage implements RepoMetaclient.
func (d metaclient) GetReposForPackage(name string) ([]v1alpha1.PackageRepository, error) {
	if repoList, err := d.clientset.client.ListPackageRepositories(context.TODO()); err != nil {
		return nil, err
	} else {
		var result []v1alpha1.PackageRepository
		var compositeErr error
		for _, repo := range repoList.Items {
			var index types.PackageRepoIndex
			if err := d.clientset.ForRepo(repo).FetchPackageRepoIndex(&index); err != nil {
				multierr.AppendInto(&compositeErr, err)
			} else {
				if slices.ContainsFunc(index.Packages, func(item types.PackageRepoIndexItem) bool { return item.Name == name }) {
					result = append(result, repo)
				}
			}
		}
		if compositeErr != nil {
			if len(result) == 0 {
				return result, compositeErr
			} else {
				return result, repoerror.Partial(compositeErr)
			}
		}
		return result, nil
	}
}
