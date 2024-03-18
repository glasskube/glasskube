package describe

import (
	"context"
	"errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/manifest"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func DescribePackage(
	ctx context.Context,
	pkgName string,
) (*v1alpha1.Package, *client.PackageStatus, *v1alpha1.PackageManifest, string, error) {
	pkgClient := client.FromContext(ctx)
	var pkg v1alpha1.Package
	var status *client.PackageStatus
	if err := pkgClient.Packages().Get(ctx, pkgName, &pkg); err == nil {
		// pkg installed: try to use installed manifest
		status = client.GetStatusOrPending(&pkg.Status)
		if installedManifest, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err == nil {
			return &pkg, status, installedManifest, "", nil
		} else if !(errors.Is(err, manifest.ErrPackageNoManifest) || apierrors.IsNotFound(err)) {
			return nil, nil, nil, "", err
		}
	} else if !apierrors.IsNotFound(err) {
		return nil, nil, nil, "", err
	}

	// pkg not installed or no manifest found: use manifest from repo
	var packageManifest v1alpha1.PackageManifest
	// TODO: Returning latestVersion in this way seems weird. We should find a better way.
	if latestVersion, err := repo.FetchLatestPackageManifest("", pkgName, &packageManifest); err != nil {
		return nil, nil, nil, "", err
	} else {
		return nil, nil, &packageManifest, latestVersion, nil
	}
}
