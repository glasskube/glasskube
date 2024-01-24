package install

import (
	"context"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// InstallBlocking creates a new v1alpha1.Package custom resource in the cluster and waits until
// the package has either status Ready or Failed.
func InstallBlocking(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, packageName string) (*client.PackageStatus, error) {
	pkg, err := Install(pkgClient, ctx, packageName)
	if err != nil {
		return nil, err
	}

	status, err := awaitInstall(pkgClient, ctx, pkg.GetUID())
	if err != nil {
		return nil, err
	}
	return status, nil
}

// Install creates a new v1alpha1.Package custom resource in the cluster.
func Install(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, packageName string) (*v1alpha1.Package, error) {
	pkg := client.NewPackage(packageName)
	err := pkgClient.Packages().Create(ctx, pkg)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Installing %v.\n", packageName)
	return pkg, err
}

func awaitInstall(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, pkgUID types.UID) (*client.PackageStatus, error) {
	watcher, err := pkgClient.Packages().Watch(ctx)
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
