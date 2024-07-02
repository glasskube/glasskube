package update

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	apitypes "k8s.io/apimachinery/pkg/types"
)

type PackagesGetter interface {
	Get(ctx context.Context) ([]ctrlpkg.Package, error)
	Explicit() bool
}

type exactPackagesGetter struct {
	pkgs []ctrlpkg.Package
}

// Explicit implements PackagesGetter.
func (s *exactPackagesGetter) Explicit() bool {
	return true
}

// Get implements PackagesGetter.
func (s *exactPackagesGetter) Get(ctx context.Context) ([]ctrlpkg.Package, error) {
	return s.pkgs, nil
}

func GetExact(pkgs []ctrlpkg.Package) PackagesGetter {
	return &exactPackagesGetter{pkgs: pkgs}
}

type clusterPackagesGetter struct{}

// Explicit implements PackagesGetter.
func (c *clusterPackagesGetter) Explicit() bool {
	return false
}

// Get implements PackagesGetter.
func (*clusterPackagesGetter) Get(ctx context.Context) ([]ctrlpkg.Package, error) {
	client := clicontext.PackageClientFromContext(ctx)
	var list v1alpha1.ClusterPackageList
	if err := client.ClusterPackages().GetAll(ctx, &list); err != nil {
		return nil, err
	} else {
		r := make([]ctrlpkg.Package, len(list.Items))
		for i := range list.Items {
			r[i] = &list.Items[i]
		}
		return r, err
	}
}

func GetAllClusterPackages() PackagesGetter {
	return &clusterPackagesGetter{}
}

type packagesGetter struct{ ns string }

// Explicit implements PackagesGetter.
func (p *packagesGetter) Explicit() bool {
	return false
}

// Get implements PackagesGetter.
func (p *packagesGetter) Get(ctx context.Context) ([]ctrlpkg.Package, error) {
	client := clicontext.PackageClientFromContext(ctx)
	var list v1alpha1.PackageList
	if err := client.Packages(p.ns).GetAll(ctx, &list); err != nil {
		return nil, err
	} else {
		r := make([]ctrlpkg.Package, len(list.Items))
		for i := range list.Items {
			r[i] = &list.Items[i]
		}
		return r, err
	}
}

func GetAllPackages(ns string) PackagesGetter {
	return &packagesGetter{ns: ns}
}

type clusterPackagesWithNameGetter struct {
	names []string
}

// Explicit implements PackagesGetter.
func (c *clusterPackagesWithNameGetter) Explicit() bool {
	return true
}

// Get implements PackagesGetter.
func (c *clusterPackagesWithNameGetter) Get(ctx context.Context) ([]ctrlpkg.Package, error) {
	client := clicontext.PackageClientFromContext(ctx)
	pkgs := make([]ctrlpkg.Package, len(c.names))
	for i, name := range c.names {
		var pkg v1alpha1.ClusterPackage
		if err := client.ClusterPackages().Get(ctx, name, &pkg); err != nil {
			return nil, err
		}
		pkgs[i] = &pkg
	}
	return pkgs, nil
}

func GetClusterPackagesWithNames(names []string) PackagesGetter {
	return &clusterPackagesWithNameGetter{names: names}
}

func GetClusterPackageWithName(name string) PackagesGetter {
	return &clusterPackagesWithNameGetter{names: []string{name}}
}

type packagesWithNameGetter struct {
	namespacedNames []apitypes.NamespacedName
}

// Explicit implements PackagesGetter.
func (p *packagesWithNameGetter) Explicit() bool {
	return true
}

// Get implements PackagesGetter.
func (p *packagesWithNameGetter) Get(ctx context.Context) ([]ctrlpkg.Package, error) {
	client := clicontext.PackageClientFromContext(ctx)
	pkgs := make([]ctrlpkg.Package, len(p.namespacedNames))
	for i, name := range p.namespacedNames {
		var pkg v1alpha1.Package
		if err := client.Packages(name.Namespace).Get(ctx, name.Name, &pkg); err != nil {
			return nil, err
		}
		pkgs[i] = &pkg
	}
	return pkgs, nil
}

func GetPackagesWithNames(namespacedNames []apitypes.NamespacedName) PackagesGetter {
	return &packagesWithNameGetter{}
}

func GetPackageWithName(namespace, name string) PackagesGetter {
	return &packagesWithNameGetter{
		namespacedNames: []apitypes.NamespacedName{
			{Namespace: namespace, Name: name},
		},
	}
}
