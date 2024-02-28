package adapter

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

type ClientAdapter interface {
	GetPackage(ctx context.Context, pkgName string) (*v1alpha1.Package, error)
}

type RepoAdapter interface {
	GetLatestVersion(repo string, pkgName string) (string, error)
	GetMaxVersionCompatibleWith(repo string, pkgName string, versionRange string) (string, error)
}
