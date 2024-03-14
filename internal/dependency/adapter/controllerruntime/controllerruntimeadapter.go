package controllerruntime

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ControllerRuntimeAdapter struct {
	client client.Client
}

func NewControllerRuntimeAdapter(client client.Client) *ControllerRuntimeAdapter {
	return &ControllerRuntimeAdapter{
		client: client,
	}
}

func (a *ControllerRuntimeAdapter) GetPackage(ctx context.Context, pkgName string) (*v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	if err := a.client.Get(ctx, types.NamespacedName{
		Name: pkgName,
	}, &pkg); err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return &pkg, nil
	}
}

func (a *ControllerRuntimeAdapter) GetPackageInfo(ctx context.Context, pkgInfoName string) (*v1alpha1.PackageInfo, error) {
	var pkgInfo v1alpha1.PackageInfo
	if err := a.client.Get(ctx, types.NamespacedName{
		Name: pkgInfoName,
	}, &pkgInfo); err != nil {
		return nil, err
	} else {
		return &pkgInfo, nil
	}
}

func (a *ControllerRuntimeAdapter) ListPackages(ctx context.Context) (*v1alpha1.PackageList, error) {
	var pkgList v1alpha1.PackageList
	if err := a.client.List(ctx, &pkgList, &client.ListOptions{}); err != nil {
		return nil, err
	}
	return &pkgList, nil
}
