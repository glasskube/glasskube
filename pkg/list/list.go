package list

import (
	"context"
	"fmt"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"sync"
)

type PackageTeaserWithStatus struct {
	PackageName      string
	ShortDescription string
	Status           *client.PackageStatus
	IconUrl          string
}

func Get(client *client.PackageV1Alpha1Client, ctx context.Context, name string) (*v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	err := client.Packages().Get(ctx, name, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func GetInstalled(client *client.PackageV1Alpha1Client, ctx context.Context) (*v1alpha1.PackageList, error) {
	ls := &v1alpha1.PackageList{}
	err := client.Packages().GetAll(ctx, ls)
	if err != nil {
		return nil, err
	}
	return ls, nil
}

func GetPackagesWithStatus(
	pkgClient *client.PackageV1Alpha1Client,
	ctx context.Context,
	onlyInstalled bool,
) ([]*PackageTeaserWithStatus, error) {
	index, installed, err := fetchRepoAndInstalled(pkgClient, ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*PackageTeaserWithStatus, 0, len(index.Packages))
	for _, description := range index.Packages {
		pkgWithStatus := &PackageTeaserWithStatus{
			PackageName:      description.Name,
			ShortDescription: description.ShortDescription,
			IconUrl:          description.IconUrl,
			Status:           nil,
		}
		for _, inst := range installed.Items {
			if description.Name == inst.Name {
				if stat := client.GetStatus(&inst.Status); stat != nil {
					pkgWithStatus.Status = stat
				} else {
					pkgWithStatus.Status = &client.PackageStatus{
						Status: "Pending",
					}
				}
				break
			}
		}
		if !onlyInstalled || pkgWithStatus.Status != nil {
			res = append(res, pkgWithStatus)
		}
	}
	return res, nil
}

func fetchRepoAndInstalled(
	pkgClient *client.PackageV1Alpha1Client,
	ctx context.Context,
) (*repo.PackageRepoIndex,
	*v1alpha1.PackageList,
	error,
) {
	var index *repo.PackageRepoIndex
	var errRepo error
	var installed *v1alpha1.PackageList
	var errInst error
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		index, errRepo = repo.FetchPackageRepoIndex(ctx, "")
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		installed, errInst = GetInstalled(pkgClient, ctx)
		wg.Done()
	}()
	wg.Wait()
	if errRepo != nil {
		return nil, nil, fmt.Errorf("could not fetch package repository index: %w\n", errRepo)
	}
	if errInst != nil {
		return nil, nil, fmt.Errorf("could not fetch installed packages: %w\n", errInst)
	}
	return index, installed, nil
}
