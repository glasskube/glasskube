package list

import (
	"context"
	"fmt"
	"sync"

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
	Status            *client.PackageStatus
	Package           *v1alpha1.Package
	InstalledManifest *v1alpha1.PackageManifest
}

type ListOptions struct {
	IncludePackageInfos bool
	OnlyInstalled       bool
	OnlyOutdated        bool
}

type lister struct {
	pkgClient  client.PackageV1Alpha1Client
	repoClient repoclient.RepoClientset
}

func NewLister(ctx context.Context) *lister {
	return &lister{
		pkgClient:  cliutils.PackageClient(ctx),
		repoClient: cliutils.RepositoryClientset(ctx),
	}
}

func (l *lister) GetPackagesWithStatus(
	ctx context.Context,
	options ListOptions,
) ([]*PackageWithStatus, error) {
	index, err := l.fetchRepoAndInstalled(ctx, options)
	result := make([]*PackageWithStatus, 0, len(index))
	for _, item := range index {
		pkgWithStatus := PackageWithStatus{
			MetaIndexItem: *item.IndexItem,
		}

		if !((options.OnlyInstalled && !item.Installed()) || (options.OnlyOutdated && !item.Outdated())) {
			pkgWithStatus.Package = item.Package
			pkgWithStatus.Status = client.GetStatusOrPending(item.Package)

			if item.PackageInfo != nil {
				pkgWithStatus.InstalledManifest = item.PackageInfo.Status.Manifest
			}

			result = append(result, &pkgWithStatus)
		}
	}
	return result, err
}

func (l *lister) fetchRepoAndInstalled(ctx context.Context, options ListOptions) (
	[]result,
	error,
) {
	var index repotypes.MetaIndex
	var packages v1alpha1.PackageList
	var packageInfos v1alpha1.PackageInfoList
	var repoErr, pkgErr, pkgInfoErr error
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := l.repoClient.Meta().FetchMetaIndex(&index); err != nil {
			repoErr = fmt.Errorf("could not fetch package repository index: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := l.pkgClient.Packages().GetAll(ctx, &packages); err != nil {
			pkgErr = fmt.Errorf("could not fetch installed packages: %w", err)
		}
	}()

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

	compositeErr := multierr.Combine(repoErr, pkgErr, pkgInfoErr)
	if pkgErr != nil || pkgInfoErr != nil {
		return nil, compositeErr
	}

	result := make([]result, len(index.Packages))
	for i, indexPackage := range index.Packages {
		result[i].IndexItem = &index.Packages[i]
		for j, clusterPackage := range packages.Items {
			if indexPackage.Name == clusterPackage.Name {
				result[i].Package = &packages.Items[j]
				packageInfoName := names.PackageInfoName(clusterPackage)
				for k, packageInfo := range packageInfos.Items {
					if packageInfo.Name == packageInfoName {
						result[i].PackageInfo = &packageInfos.Items[k]
						break
					}
				}
				break
			}
		}
	}
	return result, compositeErr
}
