package dependency

import (
	"context"
	"errors"
	"slices"
	"strings"

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

	// We can not call validate the graph here, because the initial graph, representing the current cluster state, may
	// actually be invalid. This is because we let the operator create/update dependencies which can only happen after
	// a package is already created.

	if err := dm.add(g, *manifest, version); err != nil {
		return nil, err
	}

	requirements, err := dm.addDependencies(g, manifest.Name, false)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(requirements, func(a, b Requirement) int { return strings.Compare(a.Name, b.Name) })

	var conflicts []Conflict
	for _, err := range multierr.Errors(g.Validate()) {
		if conflict, err := errorToConflict(err); err != nil {
			return nil, err
		} else {
			conflicts = append(conflicts, *conflict)
		}
	}

	status := ValidationResultStatusOk
	if len(requirements) > 0 {
		status = ValidationResultStatusResolvable
	}
	if len(conflicts) > 0 {
		status = ValidationResultStatusConflict
	}
	return &ValidationResult{
		Status:       status,
		Requirements: requirements,
		Conflicts:    conflicts,
	}, nil
}

// NewGraph constructs a DependencyGraph from all packages returned by clientAdapter.ListPackages
func (dm *DependendcyManager) NewGraph(ctx context.Context) (*graph.DependencyGraph, error) {
	pkgs, err := dm.pkgClient.ListPackages(ctx)
	if err != nil {
		return nil, err
	}
	g := graph.NewGraph()
	for _, pkg := range pkgs.Items {
		var deps []v1alpha1.Dependency
		installedVersion := pkg.Spec.PackageInfo.Version
		if !pkg.DeletionTimestamp.IsZero() {
			// A package that is currently being deleted is added to the graph, but in a state representing
			// "not installed"
			installedVersion = ""
		} else if pi, err := dm.pkgClient.GetPackageInfo(ctx, names.PackageInfoName(pkg)); err != nil {
			return nil, err
		} else if pi.Status.Manifest != nil {
			deps = pi.Status.Manifest.Dependencies
		}
		if err := g.Add(pkg.Name, installedVersion, deps, len(pkg.OwnerReferences) == 0); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (dm *DependendcyManager) add(
	g *graph.DependencyGraph,
	manifest v1alpha1.PackageManifest,
	version string,
) error {
	return g.Add(manifest.Name, version, manifest.Dependencies, g.Manual(manifest.Name))
}

// addDependencies adds the highest possible version of every uninstalled dependency and installs all transitive
// dependencies
func (dm *DependendcyManager) addDependencies(
	g *graph.DependencyGraph,
	name string,
	transitive bool,
) ([]Requirement, error) {
	var allAdded []Requirement
	for _, dep := range g.Dependencies(name) {
		if g.Version(dep) == nil {
			if versions, err := dm.getVersions(dep); err != nil {
				return nil, err
			} else if maxVersion, err := g.Max(dep, versions); err != nil {
				// This error occurs when no suitable version exists.
				// In this case, the dependency is not added to the graph and a validation error detects this later.
				continue
			} else if depManifest, err := dm.repoAdapter.GetManifest(dep, maxVersion.Original()); err != nil {
				return nil, err
			} else if err := dm.add(g, *depManifest, maxVersion.Original()); err != nil {
				return nil, err
			} else if added, err := dm.addDependencies(g, depManifest.Name, true); err != nil {
				return nil, err
			} else {
				allAdded = append(allAdded, Requirement{
					PackageWithVersion: PackageWithVersion{Name: depManifest.Name, Version: maxVersion.Original()},
					Transitive:         transitive})
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
			Actual:   PackageWithVersion{Name: errConstraint.Name, Version: version},
			Required: PackageWithVersion{Name: errConstraint.Name, Version: constraint},
			Cause:    err,
		}, nil
	} else {
		return nil, err
	}
}

// getVersions is a utility to get all versions for a package from repoAdapter and also parse them
func (dm *DependendcyManager) getVersions(name string) ([]*semver.Version, error) {
	if versions, err := dm.repoAdapter.GetVersions(name); err != nil {
		return nil, err
	} else {
		parsedVersions := make([]*semver.Version, len(versions))
		for i, version := range versions {
			parsedVersions[i], err = semver.NewVersion(version)
			if err != nil {
				return nil, err
			}
		}
		return parsedVersions, nil
	}
}
