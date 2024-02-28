package goclient

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"k8s.io/apimachinery/pkg/api/errors"
)

type GoClientAdapter struct {
	pkgClient *client.PackageV1Alpha1Client
}

func NewGoClientAdapter(pkgClient *client.PackageV1Alpha1Client) *GoClientAdapter {
	return &GoClientAdapter{
		pkgClient: pkgClient,
	}
}

func (a *GoClientAdapter) GetPackage(ctx context.Context, pkgName string) (*v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	if err := a.pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return &pkg, nil
	}
}
