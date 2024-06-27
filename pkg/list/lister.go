package list

import (
	"context"
	"fmt"
	"sync"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/names"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	repotypes "github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/pkg/client"
	"go.uber.org/multierr"
)

type PackageWithStatus struct {
	repotypes.MetaIndexItem
	Status            *client.PackageStatus     `json:"status,omitempty"`
	ClusterPackage    *v1alpha1.ClusterPackage  `json:"clusterpackage,omitempty"`
	Package           *v1alpha1.Package         `json:"package,omitempty"`
	InstalledManifest *v1alpha1.PackageManifest `json:"installedmanifest,omitempty"`
}

type PackagesWithStatus struct {
	repotypes.MetaIndexItem
	Packages []*PackageWithStatus
}

type ListOptions struct {
	IncludePackageInfos bool
	OnlyInstalled       bool
	OnlyOutdated        bool
}

type lister struct {
	pkgClient   client.PackageV1Alpha1Client
	repoClient  repoclient.RepoClientset
	useCache    bool
	cachedIndex *repotypes.MetaIndex
}

func NewLister(ctx context.Context) *lister {
	return &lister{
		pkgClient:  cliutils.PackageClient(ctx),
		repoClient: cliutils.RepositoryClientset(ctx),
	}
}

// same as NewLister, but stores the meta index after the first time it has been fetched – CAUTION: this cache
// is not concurrency-safe and therefore the listers function should only be used one at a time
func NewListerWithRepoCache(ctx context.Context) *lister {
	l := NewLister(ctx)
	l.useCache = true
	return l
}

func (l *lister) GetClusterPackagesWithStatus(
	ctx context.Context,
	options ListOptions,
) ([]*PackageWithStatus, error) {
	index, err := l.fetchRepoAndInstalled(ctx, options, includeClusterPackages)
	result := make([]*PackageWithStatus, 0, len(index))
	for _, item := range index {
		if itemShouldBeIncluded(&item, options) {
			pkgWithStatus := PackageWithStatus{
				MetaIndexItem:  *item.IndexItem,
				ClusterPackage: item.ClusterPackage,
				Status:         client.GetStatusOrPending(item.ClusterPackage),
			}
			if item.PackageInfo != nil {
				pkgWithStatus.InstalledManifest = item.PackageInfo.Status.Manifest
			}
			result = append(result, &pkgWithStatus)
		}
	}
	return result, err
}

func (l *lister) GetPackagesWithStatus(
	ctx context.Context,
	options ListOptions,
) ([]*PackagesWithStatus, error) {
	index, err := l.fetchRepoAndInstalled(ctx, options, includePackages)
	result := make([]*PackagesWithStatus, 0, len(index))
	for _, item := range index {
		// TODO itemShouldBeIncluded is wrong here – need to check outdated in the inner loop
		// if itemShouldBeIncluded(&item, options) {
		ls := make([]*PackageWithStatus, 0, len(item.Packages))
		for _, pkg := range item.Packages {
			pkgWithStatus := PackageWithStatus{
				MetaIndexItem: *item.IndexItem,
				Package:       pkg,
				Status:        client.GetStatusOrPending(pkg),
			}
			if item.PackageInfo != nil {
				pkgWithStatus.InstalledManifest = item.PackageInfo.Status.Manifest
			}
			ls = append(ls, &pkgWithStatus)
		}
		result = append(result, &PackagesWithStatus{
			MetaIndexItem: *item.IndexItem,
			Packages:      ls,
		})
		// }
	}
	return result, err
}

func itemShouldBeIncluded(item *result, options ListOptions) bool {
	return !((options.OnlyInstalled && !item.Installed()) || (options.OnlyOutdated && !item.Outdated()))
}

type typeOptions int

const (
	includeClusterPackages typeOptions = 1 << iota
	includePackages
	includeAll = includePackages | includeClusterPackages
)

func (l *lister) fetchRepoAndInstalled(ctx context.Context, options ListOptions, typeOpts typeOptions) (
	[]result,
	error,
) {
	var index repotypes.MetaIndex
	var clusterPackages v1alpha1.ClusterPackageList
	var packages v1alpha1.PackageList
	var packageInfos v1alpha1.PackageInfoList
	var repoErr, clPkgErr, pkgErr, pkgInfoErr error
	wg := new(sync.WaitGroup)

	if !l.useCache || l.cachedIndex == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.repoClient.Meta().FetchMetaIndex(&index); err != nil {
				repoErr = fmt.Errorf("could not fetch package repository index: %w", err)
			}
			l.cachedIndex = &index
		}()
	} else {
		index = *l.cachedIndex
	}

	if typeOpts&includeClusterPackages != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.pkgClient.ClusterPackages().GetAll(ctx, &clusterPackages); err != nil {
				clPkgErr = fmt.Errorf("could not fetch installed clusterpackages: %w", err)
			}
		}()
	}

	if typeOpts&includePackages != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.pkgClient.Packages("").GetAll(ctx, &packages); err != nil {
				pkgErr = fmt.Errorf("could not fetch installed packages: %w", err)
			}
		}()
	}

	if options.IncludePackageInfos {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.pkgClient.PackageInfos().GetAll(ctx, &packageInfos); err != nil {
				pkgInfoErr = fmt.Errorf("could not fetch package infos: %w", err)
			}
		}()
	}

	wg.Wait()

	compositeErr := multierr.Combine(repoErr, clPkgErr, pkgErr, pkgInfoErr)
	if clPkgErr != nil || pkgErr != nil || pkgInfoErr != nil {
		// repoErr is ignored here, because that way we can still return a partial list, if one of the repos is not available
		return nil, compositeErr
	}

	// TODO what if a package is namespaced in one repository, and with the same name cluster scoped in another??

	resultLs := make([]result, 0)
	for i, indexPackage := range index.Packages {
		res := result{
			IndexItem: &index.Packages[i],
		}
		if indexPackage.Scope.IsCluster() && typeOpts&includeClusterPackages != 0 {
			for j, pkg := range clusterPackages.Items {
				if indexPackage.Name == pkg.Name {
					res.ClusterPackage = &clusterPackages.Items[j]
					setPackageInfo(packageInfos, &res, &pkg)
					break
				}
			}
			resultLs = append(resultLs, res)
		} else if indexPackage.Scope.IsNamespaced() && typeOpts&includePackages != 0 {
			for j, pkg := range packages.Items {
				if indexPackage.Name == pkg.Spec.PackageInfo.Name {
					res.Packages = append(res.Packages, &packages.Items[j])
					setPackageInfo(packageInfos, &res, &pkg)
				}
			}
			resultLs = append(resultLs, res)
		}
	}

	return resultLs, compositeErr
}

func setPackageInfo(packageInfos v1alpha1.PackageInfoList, res *result, pkg ctrlpkg.Package) {
	packageInfoName := names.PackageInfoName(pkg)
	for k, packageInfo := range packageInfos.Items {
		if packageInfo.Name == packageInfoName {
			res.PackageInfo = &packageInfos.Items[k]
			break
		}
	}
}
