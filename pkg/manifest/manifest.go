package manifest

import (
	"context"
	"errors"

	"github.com/glasskube/glasskube/pkg/client"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

func GetInstalledManifest(ctx context.Context, pkgName string) (*v1alpha1.PackageManifest, error) {
	pkgClient := client.FromContext(ctx)
	var pkg v1alpha1.Package
	if err := pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		return nil, err
	}
	return GetInstalledManifestForPackage(ctx, pkg)
}

func GetInstalledManifestForPackage(ctx context.Context, pkg v1alpha1.Package) (*v1alpha1.PackageManifest, error) {
	pkgClient := client.FromContext(ctx)
	if len(pkg.Status.OwnedPackageInfos) == 0 {
		return nil, errors.New("Package has no owned PackageInfo")
	}
	packageInfoName := pkg.Status.OwnedPackageInfos[len(pkg.Status.OwnedPackageInfos)-1].Name
	var packageInfo v1alpha1.PackageInfo
	if err := pkgClient.PackageInfos().Get(ctx, packageInfoName, &packageInfo); err != nil {
		return nil, err
	} else if packageInfo.Status.Manifest != nil {
		return packageInfo.Status.Manifest, nil
	} else {
		return nil, errors.New("package has no manifest")
	}
}
