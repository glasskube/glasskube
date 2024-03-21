package goclient

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
)

type GoClientAdapter struct {
	pkgClient client.PackageV1Alpha1Client
}

func NewGoClientAdapter(pkgClient client.PackageV1Alpha1Client) *GoClientAdapter {
	return &GoClientAdapter{
		pkgClient: pkgClient,
	}
}

func (a *GoClientAdapter) GetPackageInfo(ctx context.Context, pkgInfoName string) (*v1alpha1.PackageInfo, error) {
	var pkgInfo v1alpha1.PackageInfo
	if err := a.pkgClient.PackageInfos().Get(ctx, pkgInfoName, &pkgInfo); err != nil {
		return nil, err
	} else {
		return &pkgInfo, nil
	}
}

func (a *GoClientAdapter) ListPackages(ctx context.Context) (*v1alpha1.PackageList, error) {
	var pkgList v1alpha1.PackageList
	if err := a.pkgClient.Packages().GetAll(ctx, &pkgList); err != nil {
		return nil, err
	}
	return &pkgList, nil
}
