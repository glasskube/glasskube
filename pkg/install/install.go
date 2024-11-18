package install

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func (obj *installer) InstallBlocking(
	ctx context.Context, pkg ctrlpkg.Package, opts metav1.CreateOptions,
) (*client.PackageStatus, error) {
	obj.status.Start()
	defer obj.status.Stop()
	pkg, err := obj.install(ctx, pkg, opts)
	if err != nil {
		return nil, err
	}
	if isDryRun(opts) {
		return &client.PackageStatus{
			Status:  string(condition.Ready),
			Reason:  "DryRun",
			Message: "Dry run - package simulated as installed and ready.",
		}, nil
	}
	return obj.awaitInstall(ctx, pkg)
}

// Install creates a new v1alpha1.ClusterPackage or v1alpha1.Package custom resource in the cluster.
func (obj *installer) Install(
	ctx context.Context, pkg ctrlpkg.Package, opts metav1.CreateOptions,
) error {
	obj.status.Start()
	defer obj.status.Stop()
	_, err := obj.install(ctx, pkg, opts)
	return err
}

func (obj *installer) install(
	ctx context.Context,
	pkg ctrlpkg.Package,
	opts metav1.CreateOptions,
) (ctrlpkg.Package, error) {
	obj.status.SetStatus(fmt.Sprintf("Installing %v...", pkg.GetName()))
	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		return pkg, obj.client.ClusterPackages().Create(ctx, pkg, opts)
	case *v1alpha1.Package:
		return pkg, obj.client.Packages(pkg.GetNamespace()).Create(ctx, pkg, opts)
	default:
		return nil, fmt.Errorf("unexpected package type: %T", pkg)
	}
}

func (obj *installer) awaitInstall(ctx context.Context, pkg ctrlpkg.Package) (*client.PackageStatus, error) {
	watcher, err := obj.watch(ctx, pkg)
	if err != nil {
		return nil, err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if obj, ok := event.Object.(ctrlpkg.Package); ok && ctrlpkg.IsSameResource(obj, pkg) {
			if event.Type == watch.Added || event.Type == watch.Modified {
				if status := client.GetStatus(obj.GetStatus()); status != nil {
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

func (i *installer) watch(ctx context.Context, pkg ctrlpkg.Package) (watch.Interface, error) {
	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		return i.client.ClusterPackages().Watch(ctx, metav1.ListOptions{})
	case *v1alpha1.Package:
		return i.client.Packages(pkg.Namespace).Watch(ctx, metav1.ListOptions{})
	default:
		return nil, fmt.Errorf("unexpected package type: %T", pkg)
	}
}
