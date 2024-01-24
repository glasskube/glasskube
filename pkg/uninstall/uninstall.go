package uninstall

import (
	"context"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
)

func Uninstall(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, pkg *v1alpha1.Package) error {
	err := pkgClient.Packages().Delete(ctx, pkg)
	if err != nil {
		return err
	}
	return nil
}
