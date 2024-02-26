package uninstall

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type uninstaller struct {
	client *client.PackageV1Alpha1Client
	status statuswriter.StatusWriter
}

var deletePropagationForeground = metav1.DeletePropagationForeground

func NewUninstaller(pkgClient *client.PackageV1Alpha1Client) *uninstaller {
	return &uninstaller{client: pkgClient, status: statuswriter.Noop()}
}

func (obj *uninstaller) WithStatusWriter(sw statuswriter.StatusWriter) *uninstaller {
	obj.status = sw
	return obj
}

// UninstallBlocking deletes the v1alpha1.Package custom resource from the
// cluster and waits until the package is fully deleted.
func (obj *uninstaller) UninstallBlocking(ctx context.Context, pkg *v1alpha1.Package) error {
	obj.status.Start()
	defer obj.status.Stop()
	pkgUID, err := obj.delete(ctx, pkg)
	if err != nil {
		return err
	}
	return obj.awaitDeletion(ctx, pkgUID)
}

// Uninstall deletes the v1alpha1.Package custom resource from the cluster.
func (obj *uninstaller) Uninstall(ctx context.Context, pkg *v1alpha1.Package) error {
	obj.status.Start()
	defer obj.status.Stop()
	_, err := obj.delete(ctx, pkg)
	return err
}

func (obj *uninstaller) delete(ctx context.Context, pkg *v1alpha1.Package) (types.UID, error) {
	obj.status.SetStatus(fmt.Sprintf("Uninstalling %v...", pkg.Name))
	err := obj.client.Packages().Delete(ctx, pkg, metav1.DeleteOptions{PropagationPolicy: &deletePropagationForeground})
	if err != nil {
		return "", err
	}
	return pkg.GetUID(), nil
}

func (obj *uninstaller) awaitDeletion(ctx context.Context, pkgUID types.UID) error {
	watcher, err := obj.client.Packages().Watch(ctx)
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if obj, ok := event.Object.(*v1alpha1.Package); ok && obj.GetUID() == pkgUID {
			if event.Type == watch.Deleted {
				return nil // Package deletion confirmed
			}
		}
	}
	return errors.New("failed to confirm package deletion")
}
