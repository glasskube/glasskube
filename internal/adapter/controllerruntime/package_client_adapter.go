package controllerruntime

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/adapter"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ControllerRuntimeAdapter struct {
	client ctrlclient.Client
}

func NewPackageClientAdapter(client ctrlclient.Client) adapter.PackageClientAdapter {
	return &ControllerRuntimeAdapter{client: client}
}

func (a *ControllerRuntimeAdapter) GetPackageInfo(ctx context.Context, pkgInfoName string) (
	*v1alpha1.PackageInfo,
	error,
) {
	var pkgInfo v1alpha1.PackageInfo
	if err := a.client.Get(ctx, types.NamespacedName{
		Name: pkgInfoName,
	}, &pkgInfo); err != nil {
		return nil, err
	} else {
		return &pkgInfo, nil
	}
}

// ListPackages implements adapter.PackageClientAdapter.
func (a *ControllerRuntimeAdapter) ListPackages(ctx context.Context, namespace string) (*v1alpha1.PackageList, error) {
	var pkgList v1alpha1.PackageList
	if err := a.client.List(ctx, &pkgList, &ctrlclient.ListOptions{Namespace: namespace}); err != nil {
		return nil, err
	}
	return &pkgList, nil
}

func (a *ControllerRuntimeAdapter) ListClusterPackages(ctx context.Context) (*v1alpha1.ClusterPackageList, error) {
	var pkgList v1alpha1.ClusterPackageList
	if err := a.client.List(ctx, &pkgList, &ctrlclient.ListOptions{}); err != nil {
		return nil, err
	}
	return &pkgList, nil
}

// GetClusterPackage implements adapter.PackageClientAdapter.
func (a *ControllerRuntimeAdapter) GetClusterPackage(ctx context.Context, name string) (*v1alpha1.ClusterPackage, error) {
	var pkg v1alpha1.ClusterPackage
	return &pkg, a.client.Get(ctx, ctrlclient.ObjectKey{Name: name}, &pkg)
}

// GetPackageRepository implements adapter.PackageClientAdapter.
func (a *ControllerRuntimeAdapter) GetPackageRepository(ctx context.Context, name string) (*v1alpha1.PackageRepository, error) {
	var repo v1alpha1.PackageRepository
	return &repo, a.client.Get(ctx, ctrlclient.ObjectKey{Name: name}, &repo)
}

// ListPackageRepositories implements adapter.PackageClientAdapter.
func (a *ControllerRuntimeAdapter) ListPackageRepositories(ctx context.Context) (*v1alpha1.PackageRepositoryList, error) {
	var list v1alpha1.PackageRepositoryList
	return &list, a.client.List(ctx, &list, &ctrlclient.ListOptions{})
}
