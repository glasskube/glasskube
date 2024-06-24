package uninstall

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type uninstaller struct {
	client client.PackageV1Alpha1Client
	status statuswriter.StatusWriter
}

var deletePropagationForeground = metav1.DeletePropagationForeground

func NewUninstaller(pkgClient client.PackageV1Alpha1Client) *uninstaller {
	return &uninstaller{client: pkgClient, status: statuswriter.Noop()}
}

func (obj *uninstaller) WithStatusWriter(sw statuswriter.StatusWriter) *uninstaller {
	obj.status = sw
	return obj
}

// UninstallBlocking deletes the v1alpha1.Package custom resource from the
// cluster and waits until the package is fully deleted.
func (obj *uninstaller) UninstallBlocking(ctx context.Context, pkg *v1alpha1.ClusterPackage) error {
	obj.status.Start()
	defer obj.status.Stop()
	err := obj.delete(ctx, pkg)
	if err != nil {
		return err
	}
	return obj.awaitDeletion(ctx, pkg.Name)
}

// Uninstall deletes the v1alpha1.Package custom resource from the cluster.
func (obj *uninstaller) Uninstall(ctx context.Context, pkg *v1alpha1.ClusterPackage) error {
	obj.status.Start()
	defer obj.status.Stop()
	return obj.delete(ctx, pkg)
}

func (obj *uninstaller) delete(ctx context.Context, pkg *v1alpha1.ClusterPackage) error {
	obj.status.SetStatus(fmt.Sprintf("Uninstalling %v...", pkg.Name))
	return obj.client.ClusterPackages().
		Delete(ctx, pkg, metav1.DeleteOptions{PropagationPolicy: &deletePropagationForeground})
}

func (obj *uninstaller) awaitDeletion(ctx context.Context, name string) error {
	watcher, err := obj.client.ClusterPackages().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if pkg, ok := event.Object.(*v1alpha1.ClusterPackage); ok && pkg.Name == name {
			if event.Type == watch.Deleted {
				return nil // Package deletion confirmed
			}
		}
	}
	return errors.New("failed to confirm package deletion")
}
