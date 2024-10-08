package dependency

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	repoerror "github.com/glasskube/glasskube/internal/repo/error"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/names"
	"go.uber.org/multierr"

	"github.com/glasskube/glasskube/internal/adapter"
	"github.com/glasskube/glasskube/internal/dependency/graph"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
)

type DependendcyManager struct {
	pkgClient   adapter.PackageClientAdapter
	repoAdapter adapter.RepoAdapter
}

func NewDependencyManager(pkgClient adapter.PackageClientAdapter, repoClient repoclient.RepoClientset) *DependendcyManager {
	return &DependendcyManager{
		pkgClient:   pkgClient,
		repoAdapter: &defaultRepoAdapter{client: repoClient},
	}
}

func (dm *DependendcyManager) Validate(
	ctx context.Context,
	name, namespace string,
	manifest *v1alpha1.PackageManifest,
	version string,
) (*ValidationResult, error) {
	if manifest == nil {
		return nil, errors.New("manifest must not be nil")
	}

	g, err := dm.NewGraph(ctx)
	if err != nil {
		return nil, err
	}

	// We do not check the validation error here, because the initial graph, representing the current cluster state, may
	// actually be invalid. This is because we let the operator create/update dependencies which can only happen after
	// a package is already created.
	// We need this error though in order to check later whether a dependency error was introduced by the action that
	// is currently validated or existed before.
	errBefore := g.Validate()

	if err := dm.add(g, name, namespace, *manifest, version); err != nil {
		return nil, err
	}

	requirements, err := dm.addDependencies(g, name, namespace, false)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(requirements, func(a, b Requirement) int { return strings.Compare(a.Name, b.Name) })

	var conflicts []Conflict
	for _, err := range multierr.Errors(g.Validate()) {
		if isErrNew(err, errBefore) {
			if conflict, err := errorToConflict(err); err != nil {
				return nil, err
			} else {
				conflicts = append(conflicts, *conflict)
			}
		}
	}

	status := ValidationResultStatusOk
	if len(requirements) > 0 {
		status = ValidationResultStatusResolvable
	}
	if len(conflicts) > 0 {
		status = ValidationResultStatusConflict
	}

	var pruned []Requirement
PruneLoop:
	for _, pkgRef := range g.Prune() {
		if pkgRef.PackageName == manifest.Name && pkgRef.Namespace == namespace && pkgRef.Name == name {
			// the currently validated package (+ its requirements) might not exist yet in the graph, and would therefore always be pruned
			continue
		}
		for _, req := range requirements {
			if req.Name == pkgRef.PackageName {
				if (req.ComponentMetadata != nil && req.ComponentMetadata.Name == pkgRef.Name &&
					req.ComponentMetadata.Namespace == pkgRef.Namespace) || req.ComponentMetadata == nil {
					continue PruneLoop
				}
			}
		}
		p := Requirement{
			PackageWithVersion: PackageWithVersion{
				Name: pkgRef.PackageName,
			},
		}
		if pkgRef.Namespace != "" {
			p.ComponentMetadata = &ComponentMetadata{
				Name:      pkgRef.Name,
				Namespace: pkgRef.Namespace,
			}
		}
		pruned = append(pruned, p)
	}
	return &ValidationResult{
		Status:       status,
		Requirements: requirements,
		Conflicts:    conflicts,
		Pruned:       pruned,
	}, nil
}

// NewGraph constructs a DependencyGraph from all packages returned by clientAdapter.ListPackages
func (dm *DependendcyManager) NewGraph(ctx context.Context) (*graph.DependencyGraph, error) {
	var allPkgs []ctrlpkg.Package
	if pkgs, err := dm.pkgClient.ListClusterPackages(ctx); err != nil {
		return nil, err
	} else {
		for i := range pkgs.Items {
			allPkgs = append(allPkgs, &pkgs.Items[i])
		}
	}

	if pkgs, err := dm.pkgClient.ListPackages(ctx, ""); err != nil {
		return nil, err
	} else {
		for i := range pkgs.Items {
			allPkgs = append(allPkgs, &pkgs.Items[i])
		}
	}

	g := graph.NewGraph()
	for _, pkg := range allPkgs {
		installedVersion := pkg.GetSpec().PackageInfo.Version
		var manifest v1alpha1.PackageManifest
		if !pkg.GetDeletionTimestamp().IsZero() {
			// A package that is currently being deleted is added to the graph, but in a state representing
			// "not installed"
			installedVersion = ""
			// we need a fake manifest with the name, in order to know the packageName of the vertex
			manifest.Name = pkg.GetSpec().PackageInfo.Name
		} else if mf, err := dm.getManifestForInstalledPkg(ctx, pkg); repoerror.IsComplete(err) {
			return nil, err
		} else {
			manifest = *mf
		}
		if pkg.IsNamespaceScoped() {
			if err := g.AddNamespaced(
				pkg.GetName(),
				pkg.GetNamespace(),
				manifest,
				installedVersion,
				!pkg.InstalledAsDependency(),
			); err != nil {
				return nil, err
			}
		} else {
			if err := g.AddCluster(manifest, installedVersion, !pkg.InstalledAsDependency()); err != nil {
				return nil, err
			}
		}
	}
	return g, nil
}

func (dm *DependendcyManager) add(
	g *graph.DependencyGraph,
	name, namespace string,
	manifest v1alpha1.PackageManifest,
	version string,
) error {
	if namespace == "" {
		return g.AddCluster(manifest, version, g.Manual(manifest.Name, ""))
	} else {
		return g.AddNamespaced(name, namespace, manifest, version, g.Manual(name, namespace))
	}
}

// addDependencies adds the highest possible version of every uninstalled dependency and installs all transitive
// dependencies
func (dm *DependendcyManager) addDependencies(
	g *graph.DependencyGraph,
	name, namespace string,
	transitive bool,
) ([]Requirement, error) {
	var allAdded []Requirement
	for _, dep := range g.Dependencies(name, namespace) {
		if g.Version(dep.Name, dep.Namespace) == nil {
			if versions, err := dm.getVersions(dep.PackageName); repoerror.IsComplete(err) {
				return nil, fmt.Errorf("failed to get version of dep package \"%v\": %w", dep.PackageName, err)
			} else if maxVersion, err := g.Max(dep.Name, dep.Namespace, versions); err != nil {
				// This error occurs when no suitable version exists.
				// In this case, the dependency is not added to the graph and a validation error detects this later.
				continue
			} else if depManifest, err := dm.repoAdapter.GetManifest(dep.PackageName, maxVersion.Original()); repoerror.IsComplete(err) {
				return nil, fmt.Errorf("failed to get manifest of dep package \"%v\" in version %v: %w", dep.PackageName, maxVersion.Original(), err)
			} else if err := dm.add(g, dep.Name, dep.Namespace, *depManifest, maxVersion.Original()); err != nil {
				return nil, err
			} else if added, err := dm.addDependencies(g, dep.Name, dep.Namespace, true); err != nil {
				return nil, err
			} else {
				req := Requirement{
					PackageWithVersion: PackageWithVersion{Name: depManifest.Name, Version: maxVersion.Original()},
					Transitive:         transitive,
				}
				if dep.Namespace != "" {
					req.ComponentMetadata = &ComponentMetadata{Name: dep.Name, Namespace: dep.Namespace}
				}
				allAdded = append(allAdded, req)
				allAdded = append(allAdded, added...)
			}
		}
	}
	return allAdded, nil
}

// errorToConflict returns a Conflict if the error is a graph.ConstraintError. Otherwise, it returns the error
// unmodified
func errorToConflict(err error) (*Conflict, error) {
	if errConstraint := (&graph.ConstraintError{}); errors.As(err, &errConstraint) {
		var version string
		var constraint string
		if errConstraint.Version != nil {
			version = errConstraint.Version.Original()
		}
		if errConstraint.Constraint != nil {
			constraint = errConstraint.Constraint.String()
		}
		return &Conflict{
			Actual:   PackageWithVersion{Name: errConstraint.Package.Name, Version: version},
			Required: PackageWithVersion{Name: errConstraint.Package.Name, Version: constraint},
			Cause:    err,
		}, nil
	} else {
		return nil, err
	}
}

func isErrNew(errCurrent error, errBefore error) bool {
	if errCurrentDep := (&graph.DependencyError{}); errors.As(errCurrent, &errCurrentDep) {
		for _, err := range multierr.Errors(errBefore) {
			if errBeforeDep := (&graph.DependencyError{}); errors.As(err, &errBeforeDep) &&
				errBeforeDep.Package == errCurrentDep.Package &&
				errBeforeDep.Dependency == errCurrentDep.Dependency {
				return false
			}

		}
	}
	return true
}

// getVersions is a utility to get all versions for a package from repoAdapter and also parse them
func (dm *DependendcyManager) getVersions(name string) ([]*semver.Version, error) {
	versions, repoErr := dm.repoAdapter.GetVersions(name)
	parsedVersions := make([]*semver.Version, len(versions))
	for i, version := range versions {
		var err error
		parsedVersions[i], err = semver.NewVersion(version)
		if err != nil {
			return nil, multierr.Append(err, repoErr)
		}
	}
	return parsedVersions, repoErr
}

func (dm *DependendcyManager) getManifestForInstalledPkg(ctx context.Context, pkg ctrlpkg.Package) (*v1alpha1.PackageManifest, error) {
	if pi, err :=
		dm.pkgClient.GetPackageInfo(ctx, names.PackageInfoName(pkg)); err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	} else if apierrors.IsNotFound(err) || (err == nil && pi.Status.Manifest == nil) {
		return dm.repoAdapter.GetManifestFromRepo(pkg.GetSpec().PackageInfo.Name, pkg.GetSpec().PackageInfo.Version, pkg.GetSpec().PackageInfo.RepositoryName)
	} else {
		return pi.Status.Manifest, nil
	}
}
