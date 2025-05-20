package uninstall

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
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
func (obj *uninstaller) UninstallBlocking(ctx context.Context, pkg ctrlpkg.Package, isDryRun bool,
	deletenamespace bool) error {
	obj.status.Start()
	defer obj.status.Stop()
	var validateDeletion bool
	var err error
	if deletenamespace {
		validateDeletion, err = validateNamespaceDeletion(ctx, pkg)
		if err != nil {
			return err
		}
	} else if !validateDeletion {
		return fmt.Errorf("Namespace deletion cannot be validated")
	}
	err = obj.delete(ctx, pkg, isDryRun)
	if err != nil {
		return err
	}
	if isDryRun {
		return nil
	}
	err = obj.awaitDeletion(ctx, pkg)
	if err != nil {
		return err
	} else {
		if validateDeletion && deletenamespace {
			err = obj.deleteNamespace(ctx, pkg.GetNamespace())
			if err != nil {
				return err
			}
			obj.status.SetStatus(fmt.Sprintf("Deleting namespace %s ...", pkg.GetNamespace()))
			err = obj.awaitNamespaceRemoval(ctx, pkg.GetNamespace())
			if err != nil {
				return err
			}

		}
	}
	return nil
}

// Uninstall deletes the v1alpha1.Package custom resource from the cluster.
func (obj *uninstaller) Uninstall(ctx context.Context, pkg ctrlpkg.Package, isDryRun bool) error {
	obj.status.Start()
	defer obj.status.Stop()
	return obj.delete(ctx, pkg, isDryRun)
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

func (obj *uninstaller) awaitNamespaceRemoval(ctx context.Context, namespace string) error {
	clientset := cliutils.KubernetesClient(ctx)
	watcher, err := clientset.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", namespace),
	})
	if err != nil {
		return err
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		if event.Type == watch.Deleted {
			return nil
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

func (uninstaller *uninstaller) deleteNamespace(ctx context.Context, namespace string) error {
	clientset := cliutils.KubernetesClient(ctx)
	err := clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("❌ error deleting namespace: %v", err)
	}
	return err
}

func validateNamespaceDeletion(ctx context.Context, pkg ctrlpkg.Package) (bool, error) {
	client := cliutils.PackageClient(ctx)
	namespace := pkg.GetNamespace()
	var pkgs v1alpha1.PackageList
	if err := client.Packages(namespace).GetAll(ctx, &pkgs); err != nil {
		return false, fmt.Errorf("❌ error listing packages in namespace: %v", err)
	}

	var namespacePackages = make([]string, 0, len(pkgs.Items))
	for _, pkg := range pkgs.Items {
		namespacePackages = append(namespacePackages, pkg.Name)
	}

	if len(namespacePackages) > 1 {
		err := fmt.Errorf("\n❌ Namespace %s cannot be deleted because it contains other packages: %v",
			namespace, strings.Join(namespacePackages, ", "))
		return false, err
	}
	return true, nil
}
