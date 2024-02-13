package describe

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/manifest"
	"k8s.io/apimachinery/pkg/api/errors"
)

func DescribePackage(
	ctx context.Context,
	pkgName string,
) (*v1alpha1.Package, *client.PackageStatus, *v1alpha1.PackageManifest, error) {
	pkgClient := client.FromContext(ctx)
	pkg, err := list.Get(pkgClient, ctx, pkgName)
	if err != nil && !errors.IsNotFound(err) {
		return nil, nil, nil, err
	}
	var status *client.PackageStatus
	var returnedManifest v1alpha1.PackageManifest
	if pkg != nil {
		// pkg installed: use installed manifest
		status = client.GetStatusOrPending(&pkg.Status)
		installedManifest, err := manifest.GetInstalledManifest(ctx, pkgName)
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
	return pkg, status, &returnedManifest, nil
}
