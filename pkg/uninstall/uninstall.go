package uninstall

import (
	"context"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/client"
)

func Uninstall(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, packageName string, force bool) (bool, error) {
	pkg := &v1alpha1.Package{}
	err := pkgClient.Packages().Get(ctx, packageName, pkg)
	if err != nil {
		return false, err
	}

	proceed := force || cliutils.YesNoPrompt(
		fmt.Sprintf("%v will be removed from your cluster. Are you sure?", packageName), false)
	if !proceed {
		return false, nil
	}

	fmt.Printf("Uninstalling %v.\n", packageName)
	err = pkgClient.Packages().Delete(ctx, pkg)
	if err != nil {
		return false, err
	}
	return true, nil
}
