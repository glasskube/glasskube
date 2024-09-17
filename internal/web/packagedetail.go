package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"

	repoerror "github.com/glasskube/glasskube/internal/repo/error"

	"go.uber.org/multierr"

	"github.com/glasskube/glasskube/internal/web/components/toast"

	"github.com/glasskube/glasskube/internal/manifestvalues"

	webutil "github.com/glasskube/glasskube/internal/web/util"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/internal/web/components/pkg_config_input"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/api/errors"
)

type packageDetailPageContext struct {
	repositoryName    string
	selectedVersion   string
	manifestName      string
	pkg               ctrlpkg.Package
	manifest          *v1alpha1.PackageManifest
	renderedComponent string
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	component := r.FormValue("component")

	var pkg *v1alpha1.Package
	var manifest *v1alpha1.PackageManifest
	if namespace != "" && name != "" && namespace != "-" && name != "-" {
		var err error
		pkg, manifest, err = describe.DescribeInstalledPackage(ctx, namespace, name)
		if err != nil && !errors.IsNotFound(err) {
			s.sendToast(w,
				toast.WithErr(fmt.Errorf("failed to fetch installed package %v/%v: %w", namespace, name, err)))
			return
		} else if errors.IsNotFound(err) {
			s.swappingRedirect(w, "/packages", "main", "main")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if pkg != nil {
			repositoryName = pkg.Spec.PackageInfo.RepositoryName
		}
	}

	s.handlePackageDetailPage(ctx, &packageDetailPageContext{
		repositoryName:    repositoryName,
		selectedVersion:   selectedVersion,
		manifestName:      manifestName,
		pkg:               pkg,
		manifest:          manifest,
		renderedComponent: component,
	}, r, w)
}

func (s *server) clusterPackageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	component := r.FormValue("component")

	pkg, manifest, err := describe.DescribeInstalledClusterPackage(ctx, pkgName)
	if err != nil && !errors.IsNotFound(err) {
		s.sendToast(w,
			toast.WithErr(fmt.Errorf("failed to fetch installed clusterpackage %v: %w", pkgName, err)))
		return
	} else if pkg != nil {
		repositoryName = pkg.Spec.PackageInfo.RepositoryName
	}

	s.handlePackageDetailPage(ctx, &packageDetailPageContext{
		repositoryName:    repositoryName,
		selectedVersion:   selectedVersion,
		manifestName:      pkgName,
		pkg:               pkg,
		manifest:          manifest,
		renderedComponent: component,
	}, r, w)
}

func (s *server) handlePackageDetailPage(ctx context.Context, d *packageDetailPageContext, r *http.Request, w http.ResponseWriter) {
	var repoErr error
	var repos []v1alpha1.PackageRepository
	var usedRepo *v1alpha1.PackageRepository
	if d.repositoryName, repos, usedRepo, repoErr = s.getRepos(
		ctx, d.manifestName, d.repositoryName); repoerror.IsComplete(repoErr) {
		s.sendToast(w, toast.WithErr(repoErr))
		return
	}

	var idx repo.PackageIndex
	var latestVersion string
	var err error
	if idx, latestVersion, d.selectedVersion, err = s.getVersions(
		d.repositoryName, d.manifestName, d.selectedVersion); err != nil {
		s.sendToast(w,
			toast.WithErr(fmt.Errorf("failed to fetch package index of %v in repo %v: %w",
				d.manifestName, d.repositoryName, multierr.Append(repoErr, err))))
		return
	}

	if d.manifest == nil {
		d.manifest = &v1alpha1.PackageManifest{}
		if err := s.repoClientset.ForRepoWithName(d.repositoryName).
			FetchPackageManifest(d.manifestName, d.selectedVersion, d.manifest); err != nil {
			s.sendToast(w,
				toast.WithErr(fmt.Errorf("failed to fetch manifest of %v (%v) in repo %v: %w", d.manifestName, d.selectedVersion, d.repositoryName, err)))
			return
		}
	}

	var validationResult *dependency.ValidationResult
	var validationErr error
	if d.pkg.IsNil() {
		if d.manifest.Scope.IsCluster() {
			validationResult, validationErr =
				s.dependencyMgr.Validate(r.Context(), d.manifestName, "", d.manifest, d.selectedVersion)
		} else {
			// In this case we don't know the actual namespace, but we can assume the default
			// TODO: make name and namespace depend on user input
			validationResult, validationErr =
				s.dependencyMgr.Validate(r.Context(), d.manifestName, d.manifest.DefaultNamespace, d.manifest, d.selectedVersion)
		}
	} else {
		validationResult, validationErr =
			s.dependencyMgr.Validate(r.Context(), d.pkg.GetName(), d.pkg.GetNamespace(), d.manifest, d.selectedVersion)
	}

	if validationErr != nil {
		s.sendToast(w,
			toast.WithErr(fmt.Errorf("failed to validate dependencies of %v (%v): %w", d.manifestName, d.selectedVersion, validationErr)))
		return
	}

	valueErrors := make(map[string]error)
	datalistOptions := make(map[string]*pkg_config_input.PkgConfigInputDatalistOptions)
	nsOptions, _ := s.getNamespaceOptions()
	if !d.pkg.IsNil() {
		pkgsOptions, _ := s.getPackagesOptions(r.Context())
		for key, v := range d.pkg.GetSpec().Values {
			if resolved, err := s.valueResolver.ResolveValue(r.Context(), v); err != nil {
				valueErrors[key] = util.GetRootCause(err)
			} else if err := manifestvalues.ValidateSingle(key, d.manifest.ValueDefinitions[key], resolved); err != nil {
				valueErrors[key] = err
			}
			if v.ValueFrom != nil {
				options, err := s.getDatalistOptions(r.Context(), v.ValueFrom, nsOptions, pkgsOptions)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
				datalistOptions[key] = options
			}
		}
	} else {
		datalistOptions[""] = &pkg_config_input.PkgConfigInputDatalistOptions{Namespaces: nsOptions}
	}

	templateData := map[string]any{
		"Package":            d.pkg,
		"Status":             client.GetStatusOrPending(d.pkg),
		"Manifest":           d.manifest,
		"LatestVersion":      latestVersion,
		"UpdateAvailable":    s.isUpdateAvailableForPkg(r.Context(), d.pkg),
		"AutoUpdate":         clientutils.AutoUpdateString(d.pkg, "Disabled"),
		"ValidationResult":   validationResult,
		"ShowConflicts":      validationResult.Status == dependency.ValidationResultStatusConflict,
		"SelectedVersion":    d.selectedVersion,
		"PackageIndex":       &idx,
		"Repositories":       repos,
		"RepositoryName":     d.repositoryName,
		"ShowConfiguration":  (!d.pkg.IsNil() && len(d.manifest.ValueDefinitions) > 0 && d.pkg.GetDeletionTimestamp().IsZero()) || d.pkg.IsNil(),
		"ValueErrors":        valueErrors,
		"DatalistOptions":    datalistOptions,
		"ShowDiscussionLink": usedRepo.IsGlasskubeRepo(),
		"PackageHref":        webutil.GetPackageHrefWithFallback(d.pkg, d.manifest),
	}

	if d.renderedComponent == "header" {
		repoErr = s.templates.pkgDetailHeaderTmpl.Execute(w, templateData)
		webutil.CheckTmplError(repoErr, fmt.Sprintf("package-detail-header (%s)", d.manifestName))
	} else {
		repoErr = s.templates.pkgPageTmpl.Execute(w, s.enrichPage(r, templateData, repoErr))
		webutil.CheckTmplError(repoErr, fmt.Sprintf("package-detail (%s)", d.manifestName))
	}
}

func (s *server) getVersions(repositoryName string, pkgName string, selectedVersion string) (repo.PackageIndex, string, string, error) {
	var idx repo.PackageIndex
	if err := s.repoClientset.ForRepoWithName(repositoryName).FetchPackageIndex(pkgName, &idx); err != nil {
		return repo.PackageIndex{}, "", "", err
	}
	latestVersion := idx.LatestVersion

	if selectedVersion == "" {
		selectedVersion = latestVersion
	} else if !slices.ContainsFunc(idx.Versions, func(item types.PackageIndexItem) bool {
		return item.Version == selectedVersion
	}) {
		selectedVersion = latestVersion
	}
	return idx, latestVersion, selectedVersion, nil
}

func (s *server) getRepos(ctx context.Context, manifestName string, repositoryName string) (
	string, []v1alpha1.PackageRepository, *v1alpha1.PackageRepository, error) {
	var repos []v1alpha1.PackageRepository
	var err error
	if repos, err = s.repoClientset.Meta().GetReposForPackage(manifestName); err != nil {
		if repoerror.IsComplete(err) {
			return "", nil, nil, err
		}
		fmt.Fprintf(os.Stderr, "error getting repos for package (but can continue): %v\n", err)
	}
	if repositoryName == "" {
		if len(repos) == 0 {
			return "", nil, nil, fmt.Errorf("%v not found in any repository", manifestName)
		}
		for _, r := range repos {
			repositoryName = r.Name
			if r.IsDefaultRepository() {
				break
			}
		}
	}

	var usedRepo v1alpha1.PackageRepository
	if err := s.pkgClient.PackageRepositories().Get(ctx, repositoryName, &usedRepo); err != nil {
		return "", nil, nil, err
	}

	return repositoryName, repos, &usedRepo, err
}
