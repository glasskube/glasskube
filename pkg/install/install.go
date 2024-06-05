package install

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type installer struct {
	client client.PackageV1Alpha1Client
	status statuswriter.StatusWriter
}

func NewInstaller(pkgClient client.PackageV1Alpha1Client) *installer {
	return &installer{client: pkgClient, status: statuswriter.Noop()}
}

func (obj *installer) WithStatusWriter(sw statuswriter.StatusWriter) *installer {
	obj.status = sw
	return obj
}

// InstallBlocking creates a new v1alpha1.Package custom resource in the cluster and waits until
// the package has either status Ready or Failed.
func (obj *installer) InstallBlocking(ctx context.Context, pkg *v1alpha1.Package,
	opts metav1.CreateOptions) (*client.PackageStatus, error) {
	obj.status.Start()
	defer obj.status.Stop()
	if isDryRun(opts) {
		_, err := obj.install(ctx, pkg, opts)
		if err != nil {
			return nil, err
		}
		return &client.PackageStatus{
			Status:  string(condition.Ready),
			Reason:  "DryRun",
			Message: "Dry run - package simulated as installed and ready.",
		}, nil
	}
	pkg, err := obj.install(ctx, pkg, opts)
	if err != nil {
		return nil, err
	}
	return obj.awaitInstall(ctx, pkg.GetUID())
}

// Install creates a new v1alpha1.Package custom resource in the cluster.
func (obj *installer) Install(ctx context.Context, pkg *v1alpha1.Package,
	opts metav1.CreateOptions) error {
	obj.status.Start()
	defer obj.status.Stop()
	_, err := obj.install(ctx, pkg, opts)
	return err
}

func (obj *installer) install(
	ctx context.Context,
	pkg *v1alpha1.Package,
	opts metav1.CreateOptions,
) (*v1alpha1.Package, error) {
	obj.status.SetStatus(fmt.Sprintf("Installing %v...", pkg.Name))
	err := obj.client.Packages().Create(ctx, pkg, opts)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (obj *installer) awaitInstall(ctx context.Context, pkgUID types.UID) (*client.PackageStatus, error) {
	watcher, err := obj.client.Packages().Watch(ctx)
	if err != nil {
		return nil, err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if obj, ok := event.Object.(*v1alpha1.Package); ok && obj.GetUID() == pkgUID {
			if event.Type == watch.Added || event.Type == watch.Modified {
				if status := client.GetStatus(&obj.Status); status != nil {
					return status, nil
				}
			} else if event.Type == watch.Deleted {
				return nil, errors.New("created package has been deleted unexpectedly")
			}
		}
	}
	return nil, errors.New("failed to confirm package installation status")
}

func isDryRun(opts metav1.CreateOptions) bool {
	for _, option := range opts.DryRun {
		if option == metav1.DryRunAll {
			return true
		}
	}
	return false
}
