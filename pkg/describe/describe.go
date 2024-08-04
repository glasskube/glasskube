package describe

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/manifest"
)

func DescribeInstalledClusterPackage(ctx context.Context, pkgName string) (
	*v1alpha1.ClusterPackage, *v1alpha1.PackageManifest, error) {
	pkgClient := cliutils.PackageClient(ctx)
	var pkg v1alpha1.ClusterPackage
	err := pkgClient.ClusterPackages().Get(ctx, pkgName, &pkg)
	if err != nil {
		return nil, nil, err
	}

	mf, err := GetManifestForPkg(ctx, &pkg)
	return &pkg, mf, err
}

func DescribeInstalledPackage(ctx context.Context, namespace string, name string) (
	*v1alpha1.Package, *v1alpha1.PackageManifest, error) {
	pkgClient := cliutils.PackageClient(ctx)
	var pkg v1alpha1.Package
	err := pkgClient.Packages(namespace).Get(ctx, name, &pkg)
	if err != nil {
		return nil, nil, err
	}

	mf, err := GetManifestForPkg(ctx, &pkg)
	return &pkg, mf, err
}

func GetManifestForPkg(ctx context.Context, pkg ctrlpkg.Package) (*v1alpha1.PackageManifest, error) {
	if installedManifest, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err == nil {
		return installedManifest, nil
	} else if !errors.Is(err, manifest.ErrPackageNoManifest) {
		return nil, err
	}

	repoClient := cliutils.RepositoryClientset(ctx)
	// pkg is installed, but has either no manifest or owned package info (yet): use manifest in this version from repo
	var packageManifest v1alpha1.PackageManifest
	err := repoClient.ForPackage(pkg).FetchPackageManifest(
		pkg.GetSpec().PackageInfo.Name, pkg.GetSpec().PackageInfo.Version, &packageManifest)
	if err != nil {
		return nil, err
	} else {
		return &packageManifest, nil
	}
}

func DescribeLatestVersion(ctx context.Context, repositoryName string, packageName string) (
	*v1alpha1.PackageManifest, string, error) {

	repoClient := cliutils.RepositoryClientset(ctx)
	if len(repositoryName) == 0 {
		if repos, err := repoClient.Meta().GetReposForPackage(packageName); err != nil {
			return nil, "", err
		} else if len(repos) == 0 {
			return nil, "", fmt.Errorf("no repo found for package %v", packageName)
		} else {
			for _, repo := range repos {
				repositoryName = repo.Name
				if !repo.IsDefaultRepository() {
					break
				}
			}
		}
	}
	var packageManifest v1alpha1.PackageManifest
	if latestVersion, err := repoClient.ForRepoWithName(repositoryName).
		FetchLatestPackageManifest(packageName, &packageManifest); err != nil {
		return nil, "", err
	} else {
		return &packageManifest, latestVersion, nil
	}
}
