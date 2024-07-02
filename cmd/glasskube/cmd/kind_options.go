package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/errors"
)

type ResourceKind string

const (
	KindUnspecified    ResourceKind = ""
	KindPackage        ResourceKind = "package"
	KindClusterPackage ResourceKind = "clusterpackage"
)

// Set implements pflag.Value.
func (kind *ResourceKind) Set(value string) error {
	switch ResourceKind(value) {
	case KindPackage, KindClusterPackage:
		*kind = ResourceKind(value)
		return nil
	default:
		return fmt.Errorf("invalid kind: %s", value)
	}
}

// String implements pflag.Value.
func (kind *ResourceKind) String() string {
	return string(*kind)
}

// Type implements pflag.Value.
func (r *ResourceKind) Type() string {
	return fmt.Sprintf("(%v|%v)", KindPackage, KindClusterPackage)
}

type KindOptions struct {
	Kind ResourceKind
}

func (opts *KindOptions) AddFlagsToCommand(cmd *cobra.Command) {
	cmd.Flags().Var(&opts.Kind, "kind", "specify the kind of the resource")
}

func DefaultKindOptions() KindOptions {
	return KindOptions{
		Kind: KindUnspecified,
	}
}

func getPackageOrClusterPackage(
	ctx context.Context, name string, kindOpts KindOptions, nsOpts NamespaceOptions) (ctrlpkg.Package, error) {

	pkgClient := cliutils.PackageClient(ctx)
	var pkg v1alpha1.Package
	var cpkg v1alpha1.ClusterPackage
	// store errors separate because multierr is not threadsafe
	var errp, errcp error
	var pkgTried, cpkgTried bool

	var wg sync.WaitGroup
	if kindOpts.Kind == KindUnspecified || kindOpts.Kind == KindPackage {
		wg.Add(1)
		pkgTried = true
		go func() {
			defer wg.Done()
			namespace := nsOpts.GetActualNamespace(ctx)
			errp = pkgClient.Packages(namespace).Get(ctx, name, &pkg)
		}()
	}
	if nsOpts.Namespace == "" &&
		(kindOpts.Kind == KindUnspecified || kindOpts.Kind == KindClusterPackage) {
		// If a namespace was specified explicitly via a flag, we don't have to try to get the ClusterPackage.
		wg.Add(1)
		cpkgTried = true
		go func() {
			defer wg.Done()
			errcp = pkgClient.ClusterPackages().Get(ctx, name, &cpkg)
		}()
	}
	wg.Wait()

	// check errors other than "not found"
	var err error
	if errp != nil && !errors.IsNotFound(errp) {
		multierr.AppendInto(&err, errp)
	}
	if errcp != nil && !errors.IsNotFound(errcp) {
		multierr.AppendInto(&err, errcp)
	}
	if err != nil {
		return nil, err
	}

	// from here on, err == nil implies "not found"
	pNotFound := !pkgTried || errp != nil
	cpNotFound := !cpkgTried || errcp != nil
	if pNotFound && cpNotFound {
		return nil, fmt.Errorf("no Package or ClusterPackage found with name %v: %w; %w", name, errp, errcp)
	} else if !pNotFound && !cpNotFound {
		return nil, fmt.Errorf("both Package and ClusterPackage found with name %v. Please specify the kind explicitly", name)
	} else if !pNotFound && cpNotFound {
		return &pkg, nil
	} else {
		return &cpkg, nil
	}
}
