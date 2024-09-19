package uninstall

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
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
func (obj *uninstaller) UninstallBlocking(ctx context.Context, pkg ctrlpkg.Package, isDryRun bool, deleteNamespace bool) error {
	if deleteNamespace {
		if pkg.IsNamespaceScoped() {
			if err := obj.isNamespaceSafeToDelete(ctx, pkg); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("cannot delete namespace for cluster-scoped package")
		}
	}
	obj.status.Start()
	defer obj.status.Stop()
	err := obj.delete(ctx, pkg, isDryRun)
	if err != nil {
		return err
	}

	if !isDryRun {
		if err := obj.awaitDeletion(ctx, pkg); err != nil {
			return err
		}
	}
	if deleteNamespace {
		return obj.deleteNamespaceBlocking(ctx, pkg.GetNamespace(), isDryRun)
	}
	return nil
}

// Uninstall deletes the v1alpha1.Package custom resource from the cluster.
func (obj *uninstaller) Uninstall(ctx context.Context, pkg ctrlpkg.Package, isDryRun bool) error {
	obj.status.Start()
	defer obj.status.Stop()
	return obj.delete(ctx, pkg, isDryRun)
}

func (obj *uninstaller) isNamespaceSafeToDelete(ctx context.Context, pkg ctrlpkg.Package) error {
	var packages v1alpha1.PackageList
	err := obj.client.Packages(pkg.GetNamespace()).GetAll(ctx, &packages)
	if err != nil {
		return err
	}
	if len(packages.Items) > 1 {
		return fmt.Errorf("namespace %s contains more than one package", pkg.GetNamespace())
	}
	return nil
}

func (uninstaller *uninstaller) delete(ctx context.Context, pkg ctrlpkg.Package, isDryRun bool) error {
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: util.Pointer(metav1.DeletePropagationForeground),
	}
	if isDryRun {
		deleteOptions.DryRun = []string{metav1.DryRunAll}
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

func (uninstaller *uninstaller) deleteNamespaceBlocking(
	ctx context.Context,
	namespace string,
	isDryRun bool,
) error {
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: util.Pointer(metav1.DeletePropagationForeground),
	}
	if isDryRun {
		deleteOptions.DryRun = []string{metav1.DryRunAll}
	}
	uninstaller.status.SetStatus(fmt.Sprintf("Deleting namespace %v...", namespace))

	clientset := clicontext.KubernetesClientFromContext(ctx)
	if err := clientset.CoreV1().Namespaces().Delete(ctx, namespace, deleteOptions); err != nil {
		return err
	}
	if !isDryRun {
		return uninstaller.awaitNamespaceDeletion(ctx, clientset, namespace)
	}
	return nil
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

func (obj *uninstaller) awaitNamespaceDeletion(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) error {
	watcher, err := clientset.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", namespace),
	})
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Deleted:
			return nil // Namespace deletion confirmed
		}
	}
	return errors.New("failed to confirm namespace deletion")
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
