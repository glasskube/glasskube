package update

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/repo"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/util/retry"
)

type UpdateTransaction struct {
	Items         []updateTransactionItem
	ConflictItems []updateTransactionItemConflict
	Requirements  []dependency.Requirement
}

func (tx UpdateTransaction) IsEmpty() bool {
	for _, item := range tx.Items {
		if item.Version != "" {
			return false
		}
	}
	return true
}

type updateTransactionItem struct {
	Package ctrlpkg.Package
	Version string
}

type updateTransactionItemConflict struct {
	updateTransactionItem
	Conflicts dependency.Conflicts
}

func (txi updateTransactionItem) UpdateRequired() bool {
	return txi.Version != ""
}

type updater struct {
	client     client.PackageV1Alpha1Client
	repoClient repoclient.RepoClientset
	status     statuswriter.StatusWriter
	dm         *dependency.DependendcyManager
}

func NewUpdater(ctx context.Context) *updater {
	return &updater{
		status:     statuswriter.Noop(),
		client:     cliutils.PackageClient(ctx),
		repoClient: cliutils.RepositoryClientset(ctx),
		dm:         cliutils.DependencyManager(ctx),
	}
}

func (c *updater) WithStatusWriter(writer statuswriter.StatusWriter) *updater {
	c.status = writer
	return c
}

func (c *updater) PrepareForVersion(
	ctx context.Context, pkg ctrlpkg.Package, pkgVersion string,
) (*UpdateTransaction, error) {
	c.status.Start()
	defer c.status.Stop()
	c.status.SetStatus("Collecting installed package")

	if !semver.IsUpgradable(pkg.GetSpec().PackageInfo.Version, pkgVersion) {
		return nil, fmt.Errorf("can't update to downgraded version or equal version")
	}

	c.status.SetStatus("Updating package index")

	var tx UpdateTransaction
	item := updateTransactionItem{Package: pkg, Version: pkgVersion}
	var manifest v1alpha1.PackageManifest
	if err := c.repoClient.ForPackage(pkg).
		FetchPackageManifest(pkg.GetSpec().PackageInfo.Name, pkgVersion, &manifest); err != nil {
		return nil, err
	} else if result, err := c.dm.Validate(ctx, &manifest, pkgVersion); err != nil {
		return nil, err
	} else if len(result.Conflicts) > 0 {
		tx.ConflictItems = append(tx.ConflictItems, updateTransactionItemConflict{item, result.Conflicts})
	} else {
		tx.Items = append(tx.Items, item)
	}

	return &tx, nil
}

func (c *updater) Prepare(ctx context.Context, getters ...PackagesGetter) (*UpdateTransaction, error) {
	c.status.Start()
	defer c.status.Stop()
	c.status.SetStatus("Collecting installed packages")
	var pkgs []ctrlpkg.Package
	var explicit bool
	for _, getter := range getters {
		explicit = explicit || getter.Explicit()
		if p, err := getter.Get(ctx); err != nil {
			return nil, err
		} else {
			pkgs = append(pkgs, p...)
		}
	}
	return c.prepare(ctx, pkgs, explicit)
}

func (c *updater) prepare(
	ctx context.Context, packagesToUpdate []ctrlpkg.Package, explicitRequest bool,
) (*UpdateTransaction, error) {
	c.status.SetStatus("Updating package index")

	requirementsSet := make(map[dependency.Requirement]struct{})
	var tx UpdateTransaction

outer:
	for _, pkg := range packagesToUpdate {
		repoClient := c.repoClient.ForPackage(pkg)
		var index repo.PackageRepoIndex
		if err := repoClient.FetchPackageRepoIndex(&index); err != nil {
			return nil, fmt.Errorf("failed to fetch index: %v", err)
		}

		for _, indexItem := range index.Packages {
			if indexItem.Name == pkg.GetSpec().PackageInfo.Name {
				if semver.IsUpgradable(pkg.GetSpec().PackageInfo.Version, indexItem.LatestVersion) {
					item := updateTransactionItem{Package: pkg, Version: indexItem.LatestVersion}
					var manifest v1alpha1.PackageManifest
					if err := repoClient.FetchPackageManifest(
						pkg.GetSpec().PackageInfo.Name, indexItem.LatestVersion, &manifest); err != nil {
						return nil, err
					}
					if result, err := c.dm.Validate(ctx, &manifest, indexItem.LatestVersion); err != nil {
						return nil, err
					} else if len(result.Conflicts) > 0 {
						// This package can't be updated due to conflicts
						tx.ConflictItems = append(tx.ConflictItems,
							updateTransactionItemConflict{item, result.Conflicts},
						)
					} else {
						for _, req := range result.Requirements {
							requirementsSet[req] = struct{}{}
						}
						// this package should be updated
						tx.Items = append(tx.Items, item)
					}
				} else if explicitRequest {
					// this package is already up-to-date but an update was requested via argument
					tx.Items = append(tx.Items, updateTransactionItem{Package: pkg})
				}
				continue outer
			}
		}
		// This can happen if a package was removed from the index for some reason.
		return nil, fmt.Errorf("package %v not found in index", pkg.GetSpec().PackageInfo.Name)
	}

	for req := range requirementsSet {
		tx.Requirements = append(tx.Requirements, req)
	}

	return &tx, nil
}

func (c *updater) Apply(ctx context.Context, tx *UpdateTransaction) ([]ctrlpkg.Package, error) {
	return c.apply(ctx, tx, false)
}

func (c *updater) ApplyBlocking(ctx context.Context, tx *UpdateTransaction) ([]ctrlpkg.Package, error) {
	return c.apply(ctx, tx, true)
}

func (c *updater) apply(ctx context.Context, tx *UpdateTransaction, blocking bool) ([]ctrlpkg.Package, error) {
	c.status.Start()
	defer c.status.Stop()
	var updatedPackages []ctrlpkg.Package
	for _, item := range tx.Items {
		if item.UpdateRequired() {
			c.status.SetStatus(fmt.Sprintf("Updating %v", item.Package.GetName()))
			err := retry.OnError(retry.DefaultRetry,
				apierrors.IsNotFound,
				func() error { return c.UpdatePackage(ctx, item.Package, item.Version) },
			)
			if err != nil {
				return nil, fmt.Errorf("could not update package %v: %w", item.Package.GetName(), err)
			}
			if blocking {
				c.status.SetStatus(fmt.Sprintf("Checking %v", item.Package.GetName()))
				if err := c.awaitUpdate(ctx, item.Package); err != nil {
					return nil, fmt.Errorf("package update for %v failed: %w", item.Package.GetName(), err)
				}
			}
			updatedPackages = append(updatedPackages, item.Package)
		}
	}
	return updatedPackages, nil
}

func (c *updater) UpdatePackage(ctx context.Context, pkg ctrlpkg.Package, version string) error {
	pkg.GetSpec().PackageInfo.Version = version
	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		return c.client.ClusterPackages().Update(ctx, pkg)
	case *v1alpha1.Package:
		return c.client.Packages(pkg.GetNamespace()).Update(ctx, pkg)
	default:
		return fmt.Errorf("unexpected object kind: %v", pkg.GroupVersionKind().Kind)
	}
}

func (c *updater) awaitUpdate(ctx context.Context, pkg ctrlpkg.Package) error {
	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		watcher, err := c.client.ClusterPackages().Watch(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		return c.await(watcher, pkg)
	case *v1alpha1.Package:
		watcher, err := c.client.Packages(pkg.Namespace).Watch(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		return c.await(watcher, pkg)
	default:
		return fmt.Errorf("unexpected object kind: %v", pkg.GroupVersionKind().Kind)
	}
}

func (c *updater) await(watcher watch.Interface, pkg ctrlpkg.Package) error {
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if eventPkg, ok := event.Object.(ctrlpkg.Package); ok && eventPkg.GetUID() == pkg.GetUID() {
			if eventPkg.GetStatus().Version == eventPkg.GetSpec().PackageInfo.Version {
				return nil
			}
			if condition := meta.FindStatusCondition(
				eventPkg.GetStatus().Conditions, string(condition.Ready)); condition != nil {
				if condition.Status == metav1.ConditionFalse {
					return fmt.Errorf("Package is not ready (reason %v): %v", condition.Reason, condition.Message)
				}
			}
		}
	}
	return errors.New("watch closed unexpectedly")
}
