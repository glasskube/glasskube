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

type PackageTeaserWithStatus struct {
	PackageName       string
	ShortDescription  string
	IconUrl           string
	Status            *client.PackageStatus
	InstalledManifest *v1alpha1.PackageManifest
}

type listOptions int64

const (
	IncludePackageInfos listOptions = 1 << iota
	OnlyInstalled
)

const (
	DefaultListOptions listOptions = 0
)

func Get(client *client.PackageV1Alpha1Client, ctx context.Context, name string) (*v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	err := client.Packages().Get(ctx, name, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func GetPackagesWithStatus(
	pkgClient *client.PackageV1Alpha1Client,
	ctx context.Context,
	options listOptions,
) ([]*PackageTeaserWithStatus, error) {
	onlyInstalled := options&OnlyInstalled != 0

	index, err := fetchRepoAndInstalled(pkgClient, ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]*PackageTeaserWithStatus, 0, len(index))
	for _, item := range index {
		pkgWithStatus := PackageTeaserWithStatus{
			PackageName:      item.Teaser.Name,
			ShortDescription: item.Teaser.ShortDescription,
			IconUrl:          item.Teaser.IconUrl,
		}
		if item.Package != nil {
			if status := client.GetStatus(&item.Package.Status); status != nil {
				pkgWithStatus.Status = status
			} else {
				pkgWithStatus.Status = &client.PackageStatus{Status: "Pending"}
			}
		}
		if item.PackageInfo != nil {
			pkgWithStatus.InstalledManifest = item.PackageInfo.Status.Manifest
		}
		if !onlyInstalled || pkgWithStatus.Status != nil {
			result = append(result, &pkgWithStatus)
		}
	}
	return result, nil
}

type listResultTuple struct {
	Teaser      *repo.PackageTeaser
	Package     *v1alpha1.Package
	PackageInfo *v1alpha1.PackageInfo
}

func fetchRepoAndInstalled(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, options listOptions) (
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

	if options&IncludePackageInfos != 0 {
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
		result[i].Teaser = &index.Packages[i]
		for j, clusterPackage := range packages.Items {
			if indexPackage.Name == clusterPackage.Name {
				result[i].Package = &packages.Items[j]
				for k, packageInfo := range packageInfos.Items {
					if packageInfo.Name == clusterPackage.Spec.PackageInfo.Name {
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
