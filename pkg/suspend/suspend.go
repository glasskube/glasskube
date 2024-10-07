package suspend

import (
	"context"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Suspend(ctx context.Context, pkg ctrlpkg.Package) (Result, error) {
	if setSuspend(pkg, true) {
		if err := doUpdate(ctx, pkg); err != nil {
			return Unknown, fmt.Errorf("suspend failed for %v %v: %w", pkg.GroupVersionKind().Kind, pkg.GetName(), err)
		}
		return Suspended, nil
	}
	return UpToDate, nil
}

func Resume(ctx context.Context, pkg ctrlpkg.Package) (Result, error) {
	if setSuspend(pkg, false) {
		if err := doUpdate(ctx, pkg); err != nil {
			return Unknown, fmt.Errorf("resume failed for %v %v: %w", pkg.GroupVersionKind().Kind, pkg.GetName(), err)
		}
		return Resumed, nil
	}
	return UpToDate, nil
}

func setSuspend(pkg ctrlpkg.Package, value bool) bool {
	if pkg.GetSpec().Suspend != value {
		pkg.GetSpec().Suspend = value
		return true
	}
	return false
}

func doUpdate(ctx context.Context, pkg ctrlpkg.Package) error {
	pkgClient := cliutils.PackageClient(ctx)
	switch p := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		return pkgClient.ClusterPackages().Update(ctx, p, metav1.UpdateOptions{})
	case *v1alpha1.Package:
		return pkgClient.Packages(pkg.GetNamespace()).Update(ctx, p, metav1.UpdateOptions{})
	default:
		return fmt.Errorf("unexpected type for ctrlpkg.Package: %T", pkg)
	}
}
