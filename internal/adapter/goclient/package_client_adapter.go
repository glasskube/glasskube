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

func (a *goClientPackageAdapter) ListPackages(ctx context.Context) (*v1alpha1.PackageList, error) {
	var pkgList v1alpha1.PackageList
	if err := a.pkgClient.Packages().GetAll(ctx, &pkgList); err != nil {
		return nil, err
	}
	return &pkgList, nil
}

// GetPackage implements adapter.PackageClientAdapter.
func (a *goClientPackageAdapter) GetPackage(ctx context.Context, name string) (*v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	return &pkg, a.pkgClient.Packages().Get(ctx, name, &pkg)
}
