package describe

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/manifest"
	"k8s.io/apimachinery/pkg/api/errors"
)

func DescribePackage(
	ctx context.Context,
	pkgName string,
) (*v1alpha1.Package, *client.PackageStatus, *v1alpha1.PackageManifest, error) {
	pkgClient := client.FromContext(ctx)
	var pkg v1alpha1.Package
	var status *client.PackageStatus
	var returnedManifest v1alpha1.PackageManifest
	if err := pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil && !errors.IsNotFound(err) {
		return nil, nil, nil, err
	} else if err == nil {
		// pkg installed: use installed manifest
		status = client.GetStatusOrPending(&pkg.Status)
		installedManifest, err := manifest.GetInstalledManifestForPackage(ctx, pkg)
		if err != nil {
			return nil, nil, nil, err
		}
		returnedManifest = *installedManifest
	} else {

		// pkg not installed: use manifest from repo
		_, err = repo.FetchLatestPackageManifest("", pkgName, &returnedManifest)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return &pkg, status, &returnedManifest, nil
}
