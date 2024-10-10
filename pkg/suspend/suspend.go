package suspend

import (
	"context"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type suspendOptions struct {
	DryRun bool
}

func (opts suspendOptions) UpdateOptions() (result metav1.UpdateOptions) {
	if opts.DryRun {
		result.DryRun = []string{metav1.DryRunAll}
	}
	return
}

type Option func(opts *suspendOptions)

func DryRun() Option {
	return func(opts *suspendOptions) { opts.DryRun = true }
}

type Options []Option

func (opts Options) Get() (result suspendOptions) {
	for _, fn := range opts {
		fn(&result)
	}
	return
}

func Suspend(ctx context.Context, pkg ctrlpkg.Package, opts ...Option) (bool, error) {
	options := Options(opts).Get()
	if setSuspend(pkg, true) {
		if err := doUpdate(ctx, pkg, options.UpdateOptions()); err != nil {
			return false, fmt.Errorf("suspend failed for %v %v: %w", pkg.GroupVersionKind().Kind, pkg.GetName(), err)
		}
		return true, nil
	}
	return false, nil
}

func Resume(ctx context.Context, pkg ctrlpkg.Package, opts ...Option) (bool, error) {
	options := Options(opts).Get()
	if setSuspend(pkg, false) {
		if err := doUpdate(ctx, pkg, options.UpdateOptions()); err != nil {
			return false, fmt.Errorf("resume failed for %v %v: %w", pkg.GroupVersionKind().Kind, pkg.GetName(), err)
		}
		return true, nil
	}
	return false, nil
}

func setSuspend(pkg ctrlpkg.Package, value bool) bool {
	if pkg.GetSpec().Suspend != value {
		pkg.GetSpec().Suspend = value
		return true
	}
	return false
}

func doUpdate(ctx context.Context, pkg ctrlpkg.Package, opts metav1.UpdateOptions) error {
	pkgClient := cliutils.PackageClient(ctx)
	switch p := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		return pkgClient.ClusterPackages().Update(ctx, p, opts)
	case *v1alpha1.Package:
		return pkgClient.Packages(pkg.GetNamespace()).Update(ctx, p, opts)
	default:
		return fmt.Errorf("unexpected type for ctrlpkg.Package: %T", pkg)
	}
}
