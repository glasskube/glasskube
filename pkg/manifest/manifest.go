package manifest

import (
	"context"

	"github.com/glasskube/glasskube/pkg/client"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

func GetInstalledManifest(ctx context.Context, pkgName string) (*v1alpha1.PackageManifest, error) {
	pkgClient := client.FromContext(ctx)
	var packageInfo v1alpha1.PackageInfo
	// TODO: Change this to use the actual package info name instead of the package name
	if err := pkgClient.PackageInfos().Get(ctx, pkgName, &packageInfo); err != nil {
		return nil, err
	} else {
		return packageInfo.Status.Manifest, nil
	}
}
