package adapter

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

type ClientAdapter interface {
	GetPackage(ctx context.Context, pkgName string) (*v1alpha1.Package, error)
	GetPackageInfo(ctx context.Context, pkgInfoName string) (*v1alpha1.PackageInfo, error)
	ListPackages(ctx context.Context) (*v1alpha1.PackageList, error)
}

type RepoAdapter interface {
	GetLatestVersion(repo string, pkgName string) (string, error)
	GetMaxVersionCompatibleWith(repo string, pkgName string, versionRange string) (string, error)
}
