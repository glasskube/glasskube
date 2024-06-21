package manifest

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/names"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

var ErrPackageNoManifest = errors.New("package has no manifest")

func GetInstalledManifest(ctx context.Context, pkgName string) (*v1alpha1.PackageManifest, error) {
	pkgClient := cliutils.PackageClient(ctx)
	var pkg v1alpha1.ClusterPackage
	if err := pkgClient.ClusterPackages().Get(ctx, pkgName, &pkg); err != nil {
		return nil, err
	}
	return GetInstalledManifestForPackage(ctx, &pkg)
}

func GetInstalledManifestForPackage(ctx context.Context, pkg ctrlpkg.Package) (*v1alpha1.PackageManifest, error) {
	pkgClient := cliutils.PackageClient(ctx)
	var packageInfo v1alpha1.PackageInfo
	if err := pkgClient.PackageInfos().Get(ctx, names.PackageInfoName(pkg), &packageInfo); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("%w: %w", ErrPackageNoManifest, err)
		} else {
			return nil, err
		}
	} else if packageInfo.Status.Manifest != nil {
		return packageInfo.Status.Manifest, nil
	} else {
		return nil, ErrPackageNoManifest
	}
}
