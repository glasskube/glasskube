package goclient

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/adapter"
	"github.com/glasskube/glasskube/pkg/client"
)

type goClientPackageAdapter struct {
	pkgClient client.PackageV1Alpha1Client
}

func NewPackageClientAdapter(pkgClient client.PackageV1Alpha1Client) adapter.PackageClientAdapter {
	return &goClientPackageAdapter{pkgClient: pkgClient}
}

func (a *goClientPackageAdapter) GetPackageInfo(ctx context.Context, pkgInfoName string) (
	*v1alpha1.PackageInfo,
	error,
) {
	var pkgInfo v1alpha1.PackageInfo
	if err := a.pkgClient.PackageInfos().Get(ctx, pkgInfoName, &pkgInfo); err != nil {
		return nil, err
	} else {
		return &pkgInfo, nil
	}
}

func (a *goClientPackageAdapter) ListClusterPackages(ctx context.Context) (*v1alpha1.ClusterPackageList, error) {
	var pkgList v1alpha1.ClusterPackageList
	if err := a.pkgClient.ClusterPackages().GetAll(ctx, &pkgList); err != nil {
		return nil, err
	}
	return &pkgList, nil
}

// GetClusterPackage implements adapter.PackageClientAdapter.
func (a *goClientPackageAdapter) GetClusterPackage(ctx context.Context, name string) (*v1alpha1.ClusterPackage, error) {
	var pkg v1alpha1.ClusterPackage
	return &pkg, a.pkgClient.ClusterPackages().Get(ctx, name, &pkg)
}

// ListPackages implements adapter.PackageClientAdapter.
func (a *goClientPackageAdapter) ListPackages(ctx context.Context, namespace string) (*v1alpha1.PackageList, error) {
	var pkgList v1alpha1.PackageList
	if err := a.pkgClient.Packages(namespace).GetAll(ctx, &pkgList); err != nil {
		return nil, err
	}
	return &pkgList, nil
}

// GetPackageRepository implements adapter.PackageClientAdapter.
func (a *goClientPackageAdapter) GetPackageRepository(
	ctx context.Context, name string) (*v1alpha1.PackageRepository, error) {
	var repo v1alpha1.PackageRepository
	return &repo, a.pkgClient.PackageRepositories().Get(ctx, name, &repo)
}

// ListPackageRepositories implements adapter.PackageClientAdapter.
func (a *goClientPackageAdapter) ListPackageRepositories(ctx context.Context) (*v1alpha1.PackageRepositoryList, error) {
	var list v1alpha1.PackageRepositoryList
	return &list, a.pkgClient.PackageRepositories().GetAll(ctx, &list)
}
