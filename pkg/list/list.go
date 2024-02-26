package list

import (
	"context"
	"fmt"
	"sync"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"go.uber.org/multierr"
)

type PackageWithStatus struct {
	repo.PackageRepoIndexItem
	Status            *client.PackageStatus
	Package           *v1alpha1.Package
	InstalledManifest *v1alpha1.PackageManifest
}

type ListOptions struct {
	IncludePackageInfos bool
	OnlyInstalled       bool
	OnlyOutdated        bool
}

func GetPackagesWithStatus(
	pkgClient *client.PackageV1Alpha1Client,
	ctx context.Context,
	options ListOptions,
) ([]*PackageWithStatus, error) {

	index, err := fetchRepoAndInstalled(pkgClient, ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]*PackageWithStatus, 0, len(index))
	for _, item := range index {
		pkgWithStatus := PackageWithStatus{
			PackageRepoIndexItem: *item.IndexItem,
		}

		if !((options.OnlyInstalled && !item.Installed()) || (options.OnlyOutdated && !item.Outdated())) {
			if item.Package != nil {
				pkgWithStatus.Package = item.Package
				pkgWithStatus.Status = client.GetStatusOrPending(&item.Package.Status)
			}

			if item.PackageInfo != nil {
				pkgWithStatus.InstalledManifest = item.PackageInfo.Status.Manifest
			}

			result = append(result, &pkgWithStatus)
		}
	}
	return result, nil
}

type listResultTuple struct {
	IndexItem   *repo.PackageRepoIndexItem
	Package     *v1alpha1.Package
	PackageInfo *v1alpha1.PackageInfo
}

func (item listResultTuple) Installed() bool {
	return item.Package != nil
}

func (item listResultTuple) Outdated() bool {
	return item.Package != nil && item.IndexItem != nil &&
		item.Package.Spec.PackageInfo.Version != "" &&
		item.Package.Spec.PackageInfo.Version != item.IndexItem.LatestVersion
}

func fetchRepoAndInstalled(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, options ListOptions) (
	[]listResultTuple,
	error,
) {
	var index repo.PackageRepoIndex
	var packages v1alpha1.PackageList
	var packageInfos v1alpha1.PackageInfoList
	var compositeErr error
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {

		if err := repo.FetchPackageRepoIndex("", &index); err != nil {
			compositeErr = multierr.Append(compositeErr, fmt.Errorf("could not fetch package repository index: %w", err))
		}
		wg.Done()
	}()

	go func() {

		if err := pkgClient.Packages().GetAll(ctx, &packages); err != nil {
			compositeErr = multierr.Append(compositeErr, fmt.Errorf("could not fetch installed packages: %w", err))
		}
		wg.Done()
	}()

	if options.IncludePackageInfos {
		wg.Add(1)
		go func() {
			if err := pkgClient.PackageInfos().GetAll(ctx, &packageInfos); err != nil {
				compositeErr = multierr.Append(compositeErr, fmt.Errorf("could not fetch package infos: %w", err))
			}
			wg.Done()
		}()
	}

	wg.Wait()
	if compositeErr != nil {
		return nil, compositeErr
	}

	result := make([]listResultTuple, len(index.Packages))
	for i, indexPackage := range index.Packages {
		result[i].IndexItem = &index.Packages[i]
		for j, clusterPackage := range packages.Items {
			if indexPackage.Name == clusterPackage.Name {
				result[i].Package = &packages.Items[j]
				for k, packageInfo := range packageInfos.Items {
					packageInfoName := clusterPackage.Status.OwnedPackageInfos[len(clusterPackage.Status.OwnedPackageInfos)-1].Name
					if packageInfo.Name == packageInfoName {
						result[i].PackageInfo = &packageInfos.Items[k]
						break
					}
				}
				break
			}
		}
	}
	return result, nil
}
