package adapter

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type PackageClientAdapter interface {
	GetPackageInfo(ctx context.Context, pkgInfoName string) (*v1alpha1.PackageInfo, error)
	ListPackages(ctx context.Context) (*v1alpha1.PackageList, error)
	GetPackage(ctx context.Context, name string) (*v1alpha1.Package, error)
}

type KubernetesClientAdapter interface {
	GetSecret(ctx context.Context, name, namespace string) (*v1.Secret, error)
	GetConfigMap(ctx context.Context, name, namespace string) (*v1.ConfigMap, error)
}

type RepoAdapter interface {
	GetVersions(repoURL string, name string) ([]string, error)
	GetManifest(repoURL string, name string, version string) (*v1alpha1.PackageManifest, error)
}
