package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"

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
	repositoryName  string
	selectedVersion string
	manifestName    string
	pkg             ctrlpkg.Package
	manifest        *v1alpha1.PackageManifest
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")

	var pkg *v1alpha1.Package
	var manifest *v1alpha1.PackageManifest
	if namespace != "" && name != "" && namespace != "-" && name != "-" {
		var err error
		pkg, manifest, err = describe.DescribeInstalledPackage(ctx, namespace, name)
		if err != nil && !errors.IsNotFound(err) {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching package details of installed package %v in namespace %v", name, namespace),
				"danger")
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
		repositoryName:  repositoryName,
		selectedVersion: selectedVersion,
		manifestName:    manifestName,
		pkg:             pkg,
		manifest:        manifest,
	}, r, w)
}

func (s *server) clusterPackageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")

	pkg, manifest, err := describe.DescribeInstalledClusterPackage(ctx, pkgName)
	if err != nil && !errors.IsNotFound(err) {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching package details of installed package %v", pkgName),
			"danger")
		return
	} else if pkg != nil {
		repositoryName = pkg.Spec.PackageInfo.RepositoryName
	}

	s.handlePackageDetailPage(ctx, &packageDetailPageContext{
		repositoryName:  repositoryName,
		selectedVersion: selectedVersion,
		manifestName:    pkgName,
		pkg:             pkg,
		manifest:        manifest,
	}, r, w)
}

func (s *server) handlePackageDetailPage(ctx context.Context, d *packageDetailPageContext, r *http.Request, w http.ResponseWriter) {
	var err error
	var repos []v1alpha1.PackageRepository
	var usedRepo *v1alpha1.PackageRepository
	if d.repositoryName, repos, usedRepo, err = s.getRepos(
		ctx, d.manifestName, d.repositoryName); err != nil {
		s.respondAlertAndLog(w, err, "", "danger")
		return
	}

	var idx repo.PackageIndex
	var latestVersion string
	if idx, latestVersion, d.selectedVersion, err = s.getVersions(
		d.repositoryName, d.manifestName, d.selectedVersion); err != nil {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching package index of %v in repository %v", d.manifestName, d.repositoryName),
			"danger")
		return
	}

	if d.manifest == nil {
		d.manifest = &v1alpha1.PackageManifest{}
		if err := s.repoClientset.ForRepoWithName(d.repositoryName).
			FetchPackageManifest(d.manifestName, d.selectedVersion, d.manifest); err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching manifest of %v in version %v in repository %v",
					d.manifestName, d.selectedVersion, d.repositoryName),
				"danger")
			return
		}
	}

	res, err := s.dependencyMgr.Validate(r.Context(), d.manifest, d.selectedVersion)
	if err != nil {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred validating dependencies of %v in version %v", d.manifestName, d.selectedVersion),
			"danger")
		return
	}

	valueErrors := make(map[string]error)
	datalistOptions := make(map[string]*pkg_config_input.PkgConfigInputDatalistOptions)
	if !d.pkg.IsNil() {
		nsOptions, _ := s.getNamespaceOptions()
		pkgsOptions, _ := s.getPackagesOptions(r.Context())
		for key, v := range d.pkg.GetSpec().Values {
			if _, err := s.valueResolver.ResolveValue(r.Context(), v); err != nil {
				valueErrors[key] = util.GetRootCause(err)
			}
			if v.ValueFrom != nil {
				options, err := s.getDatalistOptions(r.Context(), v.ValueFrom, nsOptions, pkgsOptions)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
				datalistOptions[key] = options
			}
		}
	}

	err = s.templates.pkgPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Package":            d.pkg,
		"Status":             client.GetStatusOrPending(d.pkg),
		"Manifest":           d.manifest,
		"LatestVersion":      latestVersion,
		"UpdateAvailable":    s.isUpdateAvailableForPkg(r.Context(), d.pkg),
		"AutoUpdate":         clientutils.AutoUpdateString(d.pkg, "Disabled"),
		"ValidationResult":   res,
		"ShowConflicts":      res.Status == dependency.ValidationResultStatusConflict,
		"SelectedVersion":    d.selectedVersion,
		"PackageIndex":       &idx,
		"Repositories":       repos,
		"RepositoryName":     d.repositoryName,
		"ShowConfiguration":  (!d.pkg.IsNil() && len(d.manifest.ValueDefinitions) > 0 && d.pkg.GetDeletionTimestamp().IsZero()) || d.pkg.IsNil(),
		"ValueErrors":        valueErrors,
		"DatalistOptions":    datalistOptions,
		"ShowDiscussionLink": usedRepo.IsGlasskubeRepo(),
		"PackageHref":        webutil.GetPackageHrefWithFallback(d.pkg, d.manifest),
	}, err))
	checkTmplError(err, fmt.Sprintf("package-detail (%s)", d.manifestName))
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
		fmt.Fprintf(os.Stderr, "error getting repos for package; %v", err)
	} else if repositoryName == "" {
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

	return repositoryName, repos, &usedRepo, nil
}
