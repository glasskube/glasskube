package adapter

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

type ClientAdapter interface {
	GetPackageInfo(ctx context.Context, pkgInfoName string) (*v1alpha1.PackageInfo, error)
	ListPackages(ctx context.Context) (*v1alpha1.PackageList, error)
}

type RepoAdapter interface {
	GetVersions(repoURL string, name string) ([]string, error)
	GetManifest(repoURL string, name string, version string) (*v1alpha1.PackageManifest, error)
}
