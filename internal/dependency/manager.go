package dependency

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/controller/owners/utils"

	"github.com/glasskube/glasskube/internal/dependency/adapter"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	repo2 "github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/repo/client"
)

type DependendcyManager struct {
	clientAdapter adapter.ClientAdapter
	repoAdapter   adapter.RepoAdapter
	*owners.OwnerManager
}

type ValidationResultStatus string

const (
	ValidationResultStatusOk         ValidationResultStatus = "OK"
	ValidationResultStatusResolvable ValidationResultStatus = "RESOLVABLE"
	ValidationResultStatusConflict   ValidationResultStatus = "CONFLICT"
)

type PackageWithVersion struct {
	Name    string
	Version string
}

type Requirement struct {
	PackageWithVersion
}

type Conflict struct {
	Actual   PackageWithVersion
	Required PackageWithVersion
}

func (cf Conflict) String() string {
	return fmt.Sprintf("%v (required: %v, actual: %v)", cf.Required.Name, cf.Required.Version, cf.Actual.Version)
}

type Conflicts []Conflict

func (cf Conflicts) String() string {
	s := make([]string, len(cf))
	for i, c := range cf {
		s[i] = c.String()
	}
	return strings.Join(s, ", ")
}

type ValidationResult struct {
	Status       ValidationResultStatus
	Requirements []PackageWithVersion
	Conflicts    Conflicts
}

type defaultRepoAdapter struct {
	repo client.RepoClient
}

func (a *defaultRepoAdapter) GetLatestVersion(repo string, pkgName string) (string, error) {
	return a.repo.GetLatestVersion(repo, pkgName)
}

func (a *defaultRepoAdapter) GetMaxVersionCompatibleWith(repo string, pkgName string, versionRange string) (string, error) {
	var idx repo2.PackageIndex
	if err := a.repo.FetchPackageIndex(repo, pkgName, &idx); err != nil {
		return "", err
	}
	constraint, err := semver.NewConstraint(versionRange)
	if err != nil {
		return "", err
	}
	var compatibleVersions []*semver.Version
	for _, v := range idx.Versions {
		if version, err := semver.NewVersion(v.Version); err != nil {
			continue
		} else if ok := constraint.Check(version); ok {
			compatibleVersions = append(compatibleVersions, version)
		}
	}
	if len(compatibleVersions) > 0 {
		collection := semver.Collection(compatibleVersions)
		sort.Sort(collection)
		return collection[len(collection)-1].Original(), nil
	} else {
		return "", errors.New("no compatible versions found")
	}
}

func NewDependencyManager(adapter adapter.ClientAdapter, ownerMgr *owners.OwnerManager) *DependendcyManager {
	return &DependendcyManager{
		clientAdapter: adapter,
		repoAdapter:   &defaultRepoAdapter{repo: repo2.DefaultClient},
		OwnerManager:  ownerMgr,
	}
}

func (dm *DependendcyManager) WithRepo(repo client.RepoClient) *DependendcyManager {
	dm.repoAdapter = &defaultRepoAdapter{repo: repo}
	return dm
}

func (dm *DependendcyManager) Validate(ctx context.Context, piManifest *v1alpha1.PackageManifest) (*ValidationResult, error) {
	if piManifest == nil {
		return nil, errors.New("nil not allowed")
	}
	var requirements []PackageWithVersion
	var conflicts []Conflict

	// TODO this should be parallelized (be aware that tests will fail because they expect sorted requirements/conflicts)
	for _, dependency := range piManifest.Dependencies {
		if requiredPkg, err := dm.clientAdapter.GetPackage(ctx, dependency.Name); err != nil {
			return nil, err
		} else if requiredPkg == nil {
			if req, err := dm.createRequirement(dependency); err != nil {
				return nil, err
			} else {
				requirements = append(requirements, *req)
			}
		} else if dependency.Version != "" {
			if conflict, err := dm.CheckConflict(requiredPkg.Spec.PackageInfo.Version, dependency); err != nil {
				return nil, err
			} else if conflict != nil {
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
	return &ValidationResult{
		Status:       status,
		Requirements: requirements,
		Conflicts:    conflicts,
	}, nil
}

func (dm *DependendcyManager) CheckConflict(existingVersionStr string, dependency v1alpha1.Dependency) (*Conflict, error) {
	existingVersion, err := semver.NewVersion(existingVersionStr)
	if err != nil {
		return nil, err
	}
	requiredRange, err := semver.NewConstraint(dependency.Version)
	if err != nil {
		return nil, err
	}
	if ok := requiredRange.Check(existingVersion); !ok {
		return &Conflict{
			Actual: PackageWithVersion{
				Name:    dependency.Name,
				Version: existingVersionStr,
			},
			Required: PackageWithVersion{
				Name:    dependency.Name,
				Version: dependency.Version,
			},
		}, nil
	} else {
		return nil, nil
	}
}

func (dm *DependendcyManager) createRequirement(dependency v1alpha1.Dependency) (*PackageWithVersion, error) {
	requirement := &PackageWithVersion{
		Name: dependency.Name,
	}
	if dependency.Version == "" {
		if latest, err := dm.repoAdapter.GetLatestVersion("", dependency.Name); err != nil {
			return nil, err
		} else {
			requirement.Version = latest
		}
	} else {
		if maxCompatible, err := dm.repoAdapter.GetMaxVersionCompatibleWith("", dependency.Name, dependency.Version); err != nil {
			return nil, err
		} else {
			requirement.Version = maxCompatible
		}
	}
	return requirement, nil
}

type dependentTuple struct {
	Pkg        *v1alpha1.Package
	Dependency v1alpha1.Dependency
}

func (dm *DependendcyManager) IsUpdateAllowed(ctx context.Context, pkg *v1alpha1.Package, version string) (Conflicts, error) {
	dependents, err := dm.GetDependents(ctx, pkg, GetDependentsOptions{IncludeDependencies: true})
	if err != nil {
		return nil, err
	}
	var conflicts []Conflict
	for _, dependent := range dependents {
		if dependent.Dependency.Version == "" {
			continue
		} else if c, err := dm.CheckConflict(version, dependent.Dependency); err != nil {
			return nil, err
		} else if c != nil {
			conflicts = append(conflicts, *c)
		}
	}
	return conflicts, nil
}

type GetDependentsOptions struct {
	PackageList         *v1alpha1.PackageList
	IncludeDependencies bool
}

func (dm *DependendcyManager) GetDependents(ctx context.Context, pkg *v1alpha1.Package, options GetDependentsOptions) ([]dependentTuple, error) {
	var owningPkgs []v1alpha1.Package

	if ownersOfType, err := dm.OwnersOfType(&v1alpha1.Package{}, pkg); err != nil {
		return nil, err
	} else if len(ownersOfType) > 0 {
		// if there is at least one Package owning it, means that all ownersOfType are represented as owner references
		for _, owner := range ownersOfType {
			if owningPkg, err := dm.clientAdapter.GetPackage(ctx, owner.Name); err != nil || owningPkg == nil {
				return nil, err
			} else {
				owningPkgs = append(owningPkgs, *owningPkg)
			}
		}
	} else {
		// otherwise we get the owners by iterating over all packages and search for matching OwnedPackages
		pkgRef, err := utils.ToOwnedResourceRef(dm.GetScheme(), pkg)
		if err != nil {
			return nil, err
		}

		var pkgs *v1alpha1.PackageList
		if options.PackageList != nil {
			pkgs = options.PackageList
		} else if pkgs, err = dm.clientAdapter.ListPackages(ctx); err != nil {
			return nil, err
		}

		for _, otherPkg := range pkgs.Items {
			for _, ownedPkgRef := range otherPkg.Status.OwnedPackages {
				if !ownedPkgRef.MarkedForDeletion && ownedPkgRef == pkgRef {
					owningPkgs = append(owningPkgs, otherPkg)
				}
			}
		}
	}

	var dependents []dependentTuple
	for i, owningPkg := range owningPkgs {
		if options.IncludeDependencies {
			if dt, err := dm.getDependency(ctx, &owningPkg, pkg.Name); err != nil {
				return nil, err
			} else if dt != nil {
				dependents = append(dependents, *dt)
			}
		} else {
			dependents = append(dependents, dependentTuple{Pkg: &owningPkgs[i]})
		}
	}

	return dependents, nil
}

func (dm *DependendcyManager) getDependency(ctx context.Context, pkg *v1alpha1.Package, dependency string) (*dependentTuple, error) {
	pkgInfo, err := dm.clientAdapter.GetPackageInfo(ctx, pkg.Status.OwnedPackageInfos[0].Name)
	if err != nil {
		return nil, err
	}
	if pkgInfo.Status.Manifest == nil {
		return nil, nil
	}
	for _, dep := range pkgInfo.Status.Manifest.Dependencies {
		if dep.Name == dependency {
			return &dependentTuple{
				Pkg:        pkg,
				Dependency: dep,
			}, nil
		}
	}
	return nil, nil
}
