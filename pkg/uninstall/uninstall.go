package uninstall

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type uninstaller struct {
	client client.PackageV1Alpha1Client
	status statuswriter.StatusWriter
}

func NewUninstaller(pkgClient client.PackageV1Alpha1Client) *uninstaller {
	return &uninstaller{client: pkgClient, status: statuswriter.Noop()}
}

func (obj *uninstaller) WithStatusWriter(sw statuswriter.StatusWriter) *uninstaller {
	obj.status = sw
	return obj
}

// UninstallBlocking deletes the v1alpha1.Package custom resource from the
// cluster and waits until the package is fully deleted.
func (obj *uninstaller) UninstallBlocking(ctx context.Context, pkg ctrlpkg.Package) error {
	obj.status.Start()
	defer obj.status.Stop()
	err := obj.delete(ctx, pkg)
	if err != nil {
		return err
	}
	return obj.awaitDeletion(ctx, pkg)
}

// Uninstall deletes the v1alpha1.Package custom resource from the cluster.
func (obj *uninstaller) Uninstall(ctx context.Context, pkg ctrlpkg.Package) error {
	obj.status.Start()
	defer obj.status.Stop()
	return obj.delete(ctx, pkg)
}

func (uninstaller *uninstaller) delete(ctx context.Context, pkg ctrlpkg.Package) error {
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: util.Pointer(metav1.DeletePropagationForeground),
	}

	uninstaller.status.SetStatus(fmt.Sprintf("Uninstalling %v...", pkg.GetName()))

	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		return uninstaller.client.ClusterPackages().Delete(ctx, pkg, deleteOptions)
	case *v1alpha1.Package:
		return uninstaller.client.Packages(pkg.Namespace).Delete(ctx, pkg, deleteOptions)
	default:
		return fmt.Errorf("unexpected object kind: %v", pkg.GroupVersionKind().Kind)
	}
}

func (obj *uninstaller) awaitDeletion(ctx context.Context, pkg ctrlpkg.Package) error {
	watcher, err := obj.createWatcher(ctx, pkg)
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if eventPkg, ok := event.Object.(ctrlpkg.Package); ok && ctrlpkg.IsSameResource(eventPkg, pkg) {
			if event.Type == watch.Deleted {
				return nil // Package deletion confirmed
			}
		}
	}
	return errors.New("failed to confirm package deletion")
}

func (obj *uninstaller) createWatcher(ctx context.Context, pkg ctrlpkg.Package) (watch.Interface, error) {
	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		return obj.client.ClusterPackages().Watch(ctx, metav1.ListOptions{})
	case *v1alpha1.Package:
		return obj.client.Packages(pkg.Namespace).Watch(ctx, metav1.ListOptions{})
	default:
		return nil, fmt.Errorf("unexpected object kind: %v", pkg.GroupVersionKind().Kind)
	}
}
