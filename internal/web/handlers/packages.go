package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	repotypes "github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/internal/web/types"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/update"
	"k8s.io/client-go/tools/cache"
)

type clusterPackagesTemplateData struct {
	types.TemplateContextHolder
	ClusterPackages               []*list.PackageWithStatus
	ClusterPackageUpdateAvailable map[string]bool
	UpdatesAvailable              bool
}

func GetClusterPackages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clpkgs, listErr := list.NewLister(ctx).GetClusterPackagesWithStatus(ctx, list.ListOptions{IncludePackageInfos: true})
	if listErr != nil && len(clpkgs) == 0 {
		listErr = fmt.Errorf("could not load clusterpackages: %w", listErr)
		fmt.Fprintf(os.Stderr, "%v\n", listErr)
	}

	// Call isUpdateAvailable for each installed clusterpackage.
	// This is not the same as getting all updates in a single transaction, because some dependency
	// conflicts could be resolvable by installing individual clpkgs.
	installedClpkgs := make([]ctrlpkg.Package, 0, len(clpkgs))
	clpkgUpdateAvailable := map[string]bool{}
	for _, pkg := range clpkgs {
		if pkg.ClusterPackage != nil {
			installedClpkgs = append(installedClpkgs, pkg.ClusterPackage)
		}
		clpkgUpdateAvailable[pkg.Name] = isUpdateAvailableForPkg(ctx, pkg.ClusterPackage)
	}

	overallUpdatesAvailable := isUpdateAvailable(ctx, installedClpkgs)

	responder.SendPage(w, r, "pages/clusterpackages",
		responder.ContextualizedTemplate(&clusterPackagesTemplateData{
			ClusterPackages:               clpkgs,
			ClusterPackageUpdateAvailable: clpkgUpdateAvailable,
			UpdatesAvailable:              overallUpdatesAvailable,
		}),
		responder.WithPartialErr(listErr))
}

type packagesTemplateData struct {
	types.TemplateContextHolder
	InstalledPackages      []*list.PackagesWithStatus
	AvailablePackages      []*repotypes.PackageRepoIndexItem
	PackageUpdateAvailable map[string]bool
	UpdatesAvailable       bool
}

func GetPackages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allPkgs, listErr := list.NewLister(ctx).GetPackagesWithStatus(ctx, list.ListOptions{IncludePackageInfos: true})
	if listErr != nil {
		listErr = fmt.Errorf("could not load packages: %w", listErr)
		fmt.Fprintf(os.Stderr, "%v\n", listErr)
	}

	packageUpdateAvailable := map[string]bool{}
	var installed []*list.PackagesWithStatus
	var available []*repotypes.PackageRepoIndexItem
	var installedPkgs []ctrlpkg.Package
	for _, pkgsWithStatus := range allPkgs {
		if len(pkgsWithStatus.Packages) > 0 {
			for _, pkgWithStatus := range pkgsWithStatus.Packages {
				installedPkgs = append(installedPkgs, pkgWithStatus.Package)

				// Call isUpdateAvailable for each installed package.
				// This is not the same as getting all updates in a single transaction, because some dependency
				// conflicts could be resolvable by installing individual packages.
				packageUpdateAvailable[cache.MetaObjectToName(pkgWithStatus.Package).String()] =
					isUpdateAvailableForPkg(ctx, pkgWithStatus.Package)
			}
			installed = append(installed, pkgsWithStatus)
		} else {
			available = append(available, &pkgsWithStatus.PackageRepoIndexItem)
		}
	}

	overallUpdatesAvailable := false
	if len(installedPkgs) > 0 {
		overallUpdatesAvailable = isUpdateAvailable(ctx, installedPkgs)
	}

	responder.SendPage(w, r, "pages/packages",
		responder.ContextualizedTemplate(&packagesTemplateData{
			InstalledPackages:      installed,
			AvailablePackages:      available,
			PackageUpdateAvailable: packageUpdateAvailable,
			UpdatesAvailable:       overallUpdatesAvailable,
		}),
		responder.WithPartialErr(listErr))
}

func isUpdateAvailableForPkg(ctx context.Context, pkg ctrlpkg.Package) bool {
	if pkg.IsNil() {
		return false
	}
	return isUpdateAvailable(ctx, []ctrlpkg.Package{pkg})
}

func isUpdateAvailable(ctx context.Context, pkgs []ctrlpkg.Package) bool {
	if tx, err := update.NewUpdater(ctx).Prepare(ctx, update.GetExact(pkgs)); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		return false
	} else if len(tx.ConflictItems) > 0 {
		return true
	} else {
		return !tx.IsEmpty()
	}
}
