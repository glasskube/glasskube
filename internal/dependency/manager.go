package dependency

import (
	"context"
	"errors"
	"sort"

	"github.com/glasskube/glasskube/internal/dependency/adapter"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	repo2 "github.com/glasskube/glasskube/internal/repo"
)

type DependendcyManager struct {
	clientAdapter adapter.ClientAdapter
	repoAdapter   adapter.RepoAdapter
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

type ValidationResult struct {
	Status       ValidationResultStatus
	Requirements []PackageWithVersion
	Conflicts    []Conflict
}

type defaultRepoAdapter struct {
}

func (a *defaultRepoAdapter) GetLatestVersion(repo string, pkgName string) (string, error) {
	return repo2.GetLatestVersion(repo, pkgName)
}

func (a *defaultRepoAdapter) GetMaxVersionCompatibleWith(repo string, pkgName string, versionRange string) (string, error) {
	var idx repo2.PackageIndex
	if err := repo2.FetchPackageIndex(repo, pkgName, &idx); err != nil {
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

func NewDependencyManager(adapter adapter.ClientAdapter) *DependendcyManager {
	return &DependendcyManager{
		clientAdapter: adapter,
		repoAdapter:   &defaultRepoAdapter{},
	}
}

func (dm *DependendcyManager) Validate(ctx context.Context, pkg *v1alpha1.Package, pkgInfo *v1alpha1.PackageInfo) (*ValidationResult, error) {
	if pkg == nil || pkgInfo == nil || pkgInfo.Status.Manifest == nil {
		return nil, errors.New("nil not allowed")
	}
	piManifest := pkgInfo.Status.Manifest
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
			if conflict, err := dm.checkConflict(requiredPkg.Spec.PackageInfo.Version, dependency); err != nil {
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

func (dm *DependendcyManager) checkConflict(existingVersionStr string, dependency v1alpha1.Dependency) (*Conflict, error) {
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
