package update

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UpdateTransaction struct {
	Items []updateTransactionItem
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
	Package v1alpha1.Package
	Version string
}

func (txi updateTransactionItem) UpdateRequired() bool {
	return txi.Version != ""
}

type updater struct {
	client *client.PackageV1Alpha1Client
	status statuswriter.StatusWriter
}

func NewUpdater(client *client.PackageV1Alpha1Client) *updater {
	return &updater{
		client: client,
	}
}

func (c *updater) WithStatusWriter(writer statuswriter.StatusWriter) *updater {
	c.status = writer
	return c
}

func (c *updater) Prepare(ctx context.Context, packageNames []string) (*UpdateTransaction, error) {
	c.status.Start()
	defer c.status.Stop()
	c.status.SetStatus("Collecting installed packages")
	var packagesToUpdate []v1alpha1.Package
	if len(packageNames) > 0 {
		// Fetch all requested packages individually.
		// This way, we can fail early if a requested package is not installed.
		for _, name := range packageNames {
			var pkg v1alpha1.Package
			if err := c.client.Packages().Get(ctx, name, &pkg); err != nil {
				return nil, fmt.Errorf("failed to get package %v: %v", name, err)
			}
			packagesToUpdate = append(packagesToUpdate, pkg)
		}
	} else {
		var packageList v1alpha1.PackageList
		if err := c.client.Packages().GetAll(ctx, &packageList); err != nil {
			return nil, fmt.Errorf("failed to get list of installed packages: %v", err)
		}
		packagesToUpdate = packageList.Items
	}

	c.status.SetStatus("Updating package index")
	var index repo.PackageRepoIndex
	if err := repo.FetchPackageRepoIndex("", &index); err != nil {
		return nil, fmt.Errorf("failed to fetch index: %v", err)
	}

	var tx UpdateTransaction
outer:
	for _, pkg := range packagesToUpdate {
		for _, indexItem := range index.Packages {
			if indexItem.Name == pkg.Name {
				if pkg.Spec.PackageInfo.Version != "" && pkg.Spec.PackageInfo.Version != indexItem.LatestVersion {
					// this package should be updated
					tx.Items = append(tx.Items, updateTransactionItem{
						Package: pkg,
						Version: indexItem.LatestVersion,
					})
				} else if len(packageNames) > 0 {
					// this package is already up-to-date but an update was requested via argument
					tx.Items = append(tx.Items, updateTransactionItem{Package: pkg})
				}
				continue outer
			}
		}
		// This can happen if a package was removed from the index for some reason.
		return nil, fmt.Errorf("package %v not found in index", pkg.Name)
	}

	return &tx, nil
}

func (c *updater) Apply(ctx context.Context, tx *UpdateTransaction) error {
	c.status.Start()
	defer c.status.Stop()
	for _, item := range tx.Items {
		if item.UpdateRequired() {
			c.status.SetStatus(fmt.Sprintf("Updating %v", item.Package.Name))
			if err := c.UpdatePackage(ctx, &item.Package, item.Version); err != nil {
				return fmt.Errorf("could not update package %v: %w", item.Package.Name, err)
			}
			c.status.SetStatus(fmt.Sprintf("Checking %v", item.Package.Name))
			if err := c.awaitUpdate(ctx, &item.Package); err != nil {
				return fmt.Errorf("package update for %v failed: %w", item.Package.Name, err)
			}
		}
	}
	return nil
}

func (c *updater) UpdatePackage(ctx context.Context, pkg *v1alpha1.Package, version string) error {
	pkg.Spec.PackageInfo.Version = version
	return c.client.Packages().Update(ctx, pkg)
}

func (c *updater) awaitUpdate(ctx context.Context, pkg *v1alpha1.Package) error {
	watcher, err := c.client.Packages().Watch(ctx)
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if eventPkg, ok := event.Object.(*v1alpha1.Package); ok && eventPkg.GetUID() == pkg.GetUID() {
			if eventPkg.Status.Version == eventPkg.Spec.PackageInfo.Version {
				return nil
			}
			if condition := meta.FindStatusCondition(eventPkg.Status.Conditions, string(condition.Ready)); condition != nil {
				if condition.Status == metav1.ConditionFalse {
					return fmt.Errorf("Package is not ready (reason %v): %v", condition.Reason, condition.Message)
				}
			}
		}
	}
	return errors.New("watch closed unexpectedly")
}
