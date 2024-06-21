package adapter

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type PackageClientAdapter interface {
	GetPackageInfo(ctx context.Context, pkgInfoName string) (*v1alpha1.PackageInfo, error)
	ListClusterPackages(ctx context.Context) (*v1alpha1.ClusterPackageList, error)
	GetClusterPackage(ctx context.Context, name string) (*v1alpha1.ClusterPackage, error)
	ListPackages(ctx context.Context, namespace string) (*v1alpha1.PackageList, error)
	ListPackageRepositories(ctx context.Context) (*v1alpha1.PackageRepositoryList, error)
	GetPackageRepository(ctx context.Context, name string) (*v1alpha1.PackageRepository, error)
}

type KubernetesClientAdapter interface {
	GetSecret(ctx context.Context, name, namespace string) (*v1.Secret, error)
	GetConfigMap(ctx context.Context, name, namespace string) (*v1.ConfigMap, error)
}

type RepoAdapter interface {
	GetVersions(name string) ([]string, error)
	GetManifest(name string, version string) (*v1alpha1.PackageManifest, error)
}
