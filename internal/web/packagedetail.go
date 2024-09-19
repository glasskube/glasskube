package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/glasskube/glasskube/internal/namespaces"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/manifest"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

type packageContextRequest struct {
	repositoryName string
	version        string
	manifestName   string
	namespace      string
	name           string
	component      string
}

func (r packageContextRequest) namespaceAndNameSet() bool {
	return r.namespace != "" && r.name != "" && r.namespace != "-" && r.name != "-"
}

type packageContext struct {
	request  packageContextRequest
	pkg      ctrlpkg.Package
	manifest *v1alpha1.PackageManifest
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := packageContext{
		request: packageContextRequest{
			manifestName:   mux.Vars(r)["manifestName"],
			namespace:      mux.Vars(r)["namespace"],
			name:           mux.Vars(r)["name"],
			repositoryName: r.FormValue("repositoryName"),
			version:        r.FormValue("version"),
			component:      r.FormValue("component"),
		},
	}

	switch r.Method {
	case http.MethodPost:
		s.installOrConfigurePackage(w, r, &p.request)
	case http.MethodGet:
		var pkg *v1alpha1.Package
		var mf *v1alpha1.PackageManifest
		if p.request.namespaceAndNameSet() {
			var err error
			pkg, mf, err = describe.DescribeInstalledPackage(ctx, p.request.namespace, p.request.name)
			if err != nil && !errors.IsNotFound(err) {
				s.sendToast(w,
					toast.WithErr(fmt.Errorf("failed to fetch installed package %v/%v: %w", p.request.namespace, p.request.name, err)))
				return
			} else if errors.IsNotFound(err) {
				s.swappingRedirect(w, "/packages", "main", "main")
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
		p.pkg = pkg
		p.manifest = mf
		s.renderPackageDetailPage(ctx, r, w, &p)
	}
}

func (s *server) clusterPackageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := packageContext{
		request: packageContextRequest{
			manifestName:   mux.Vars(r)["pkgName"],
			repositoryName: r.FormValue("repositoryName"),
			version:        r.FormValue("version"),
			component:      r.FormValue("component"),
		},
	}

	if r.Method == http.MethodPost {
		s.installOrConfigureClusterPackage(w, r, &p.request)
	} else if r.Method == http.MethodGet {
		pkg, mf, err := describe.DescribeInstalledClusterPackage(ctx, p.request.manifestName)
		if err != nil && !errors.IsNotFound(err) {
			s.sendToast(w,
				toast.WithErr(fmt.Errorf("failed to fetch installed clusterpackage %v: %w", p.request.manifestName, err)))
			return
		}
		p.pkg = pkg
		p.manifest = mf
		s.renderPackageDetailPage(ctx, r, w, &p)
	}
}

func (s *server) renderPackageDetailPage(ctx context.Context, r *http.Request, w http.ResponseWriter, p *packageContext) {
	migrateManifest := false
	if !p.pkg.IsNil() {
		// for installed packages, the installed repo + version is the fallback, if they are not requested explicitly
		if p.request.repositoryName == "" {
			p.request.repositoryName = p.pkg.GetSpec().PackageInfo.RepositoryName
		}
		if p.request.version == "" {
			p.request.version = p.pkg.GetSpec().PackageInfo.Version
		}
		migrateManifest =
			p.request.repositoryName != p.pkg.GetSpec().PackageInfo.RepositoryName ||
				p.request.version != p.pkg.GetSpec().PackageInfo.Version
	}

	var repoErr error
	var repos []v1alpha1.PackageRepository
	var usedRepo *v1alpha1.PackageRepository
	if p.request.repositoryName, repos, usedRepo, repoErr = s.resolveRepos(
		ctx, p.request.manifestName, p.request.repositoryName); repoerror.IsComplete(repoErr) {
		s.sendToast(w, toast.WithErr(repoErr))
		return
	}

	var idx repo.PackageIndex
	var latestVersion string
	var err error
	if idx, latestVersion, p.request.version, err = s.resolveVersions(
		p.request.repositoryName, p.request.manifestName, p.request.version); err != nil {
		s.sendToast(w,
			toast.WithErr(fmt.Errorf("failed to fetch package index of %v in repo %v: %w",
				p.request.manifestName, p.request.repositoryName, multierr.Append(repoErr, err))))
		return
	}

	if p.manifest == nil || migrateManifest {
		p.manifest = &v1alpha1.PackageManifest{}
		if err := s.repoClientset.ForRepoWithName(p.request.repositoryName).
			FetchPackageManifest(p.request.manifestName, p.request.version, p.manifest); err != nil {
			s.sendToast(w,
				toast.WithErr(fmt.Errorf("failed to fetch manifest of %v (%v) in repo %v: %w",
					p.request.manifestName, p.request.version, p.request.repositoryName, err)))
			return
		}
	}

	// TODO check if dependency validation is even needed if the package is already installed??
	// probably not if the repo/version did not change, also not if only header component is rendered
	var validationResult *dependency.ValidationResult
	var validationErr error
	if p.pkg.IsNil() {
		if p.manifest.Scope.IsCluster() {
			validationResult, validationErr =
				s.dependencyMgr.Validate(r.Context(), p.request.manifestName, "", p.manifest, p.request.version)
		} else {
			// In this case we don't know the actual namespace, but we can assume the default
			// TODO: make name and namespace depend on user input
			validationResult, validationErr =
				s.dependencyMgr.Validate(r.Context(), p.request.manifestName, p.manifest.DefaultNamespace, p.manifest, p.request.version)
		}
	} else {
		validationResult, validationErr =
			s.dependencyMgr.Validate(r.Context(), p.pkg.GetName(), p.pkg.GetNamespace(), p.manifest, p.request.version)
	}

	if validationErr != nil {
		s.sendToast(w,
			toast.WithErr(fmt.Errorf("failed to validate dependencies of %v (%v): %w",
				p.request.manifestName, p.request.version, validationErr)))
		return
	}

	var lostValueDefinitions []string
	valueErrors := make(map[string]error)
	datalistOptions := make(map[string]*pkg_config_input.PkgConfigInputDatalistOptions)
	nsOptions, _ := s.getNamespaceOptions()
	if !p.pkg.IsNil() {
		pkgsOptions, _ := s.getPackagesOptions(r.Context())
		for key, v := range p.pkg.GetSpec().Values {
			if resolved, err := s.valueResolver.ResolveValue(r.Context(), v); err != nil {
				valueErrors[key] = util.GetRootCause(err)
			} else if valDef, exists := p.manifest.ValueDefinitions[key]; !exists {
				// can happen when a different repo/version of the pkg is requested (advanced options)
				lostValueDefinitions = append(lostValueDefinitions, key)
			} else if err := manifestvalues.ValidateSingle(key, valDef, resolved); err != nil {
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
	}
	datalistOptions[""] = &pkg_config_input.PkgConfigInputDatalistOptions{Namespaces: nsOptions}

	advancedOptions, err := getAdvancedOptionsFromCookie(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get advanced options from cookie: %v\n", err)
	}
	templateData := map[string]any{
		"Package":              p.pkg,
		"Status":               client.GetStatusOrPending(p.pkg),
		"Manifest":             p.manifest,
		"LatestVersion":        latestVersion,
		"UpdateAvailable":      s.isUpdateAvailableForPkg(r.Context(), p.pkg),
		"AutoUpdate":           clientutils.AutoUpdateString(p.pkg, "Disabled"),
		"ValidationResult":     validationResult,
		"ShowConflicts":        validationResult.Status == dependency.ValidationResultStatusConflict,
		"SelectedVersion":      p.request.version,
		"PackageIndex":         &idx,
		"Repositories":         repos,
		"RepositoryName":       p.request.repositoryName,
		"ShowConfiguration":    (!p.pkg.IsNil() && len(p.manifest.ValueDefinitions) > 0 && p.pkg.GetDeletionTimestamp().IsZero()) || p.pkg.IsNil(),
		"ValueErrors":          valueErrors,
		"DatalistOptions":      datalistOptions,
		"ShowDiscussionLink":   usedRepo.IsGlasskubeRepo(),
		"PackageHref":          webutil.GetPackageHrefWithFallback(p.pkg, p.manifest),
		"AdvancedOptions":      advancedOptions,
		"LostValueDefinitions": lostValueDefinitions,
	}

	if p.request.component == "header" {
		repoErr = s.templates.pkgDetailHeaderTmpl.Execute(w, templateData)
		webutil.CheckTmplError(repoErr, fmt.Sprintf("package-detail-header (%s)", p.request.manifestName))
	} else {
		repoErr = s.templates.pkgPageTmpl.Execute(w, s.enrichPage(r, templateData, repoErr))
		webutil.CheckTmplError(repoErr, fmt.Sprintf("package-detail (%s)", p.request.manifestName))
	}
}

func (s *server) resolveVersions(repositoryName string, pkgName string, selectedVersion string) (repo.PackageIndex, string, string, error) {
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

func (s *server) resolveRepos(ctx context.Context, manifestName string, repositoryName string) (
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

// installOrConfigurePackage is like installOrConfigureClusterPackage but for packages
func (s *server) installOrConfigurePackage(w http.ResponseWriter, r *http.Request, p *packageContextRequest) {
	ctx := r.Context()
	namespace := r.FormValue("namespace")
	name := r.FormValue("name")
	enableAutoUpdate := r.FormValue("enableAutoUpdate")
	dryRun, _ := strconv.ParseBool(r.FormValue("dryRun"))

	var err error
	pkg := &v1alpha1.Package{}
	var mf *v1alpha1.PackageManifest
	if err := s.pkgClient.Packages(p.namespace).Get(ctx, p.name, pkg); err != nil && !errors.IsNotFound(err) {
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to fetch package %v/%v: %w", p.namespace, p.name, err)))
		return
	} else if err != nil {
		pkg = nil
	} else {
		// because disabled form elements are not submitted, we need to fall back on repo and version if advanced options disabled
		if p.repositoryName == "" {
			p.repositoryName = pkg.Spec.PackageInfo.RepositoryName
		}
		if p.version == "" {
			p.version = pkg.Spec.PackageInfo.Version
		}
	}

	mf, err = s.resolveManifest(ctx, pkg, p.repositoryName, p.manifestName, p.version)
	if repoerror.IsPartial(err) {
		fmt.Fprintf(os.Stderr, "problem fetching manifest and repo, but installation can continue: %v", err)
	} else if err != nil {
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to get manifest and repo of %v: %w", p.manifestName, err)))
		return
	}

	if values, err := extractValues(r, mf); err != nil {
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to parse values: %w", err)))
		return
	} else if pkg == nil {
		opts := v1.CreateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		if exists, err := namespaces.Exists(ctx, s.k8sClient, namespace); err != nil {
			s.sendToast(w, toast.WithErr(fmt.Errorf("failed to check namespace: %w", err)))
			return
		} else if !exists {
			ns := v12.Namespace{
				ObjectMeta: v1.ObjectMeta{
					Name: namespace,
				},
			}
			if _, err := s.k8sClient.CoreV1().Namespaces().Create(ctx, &ns, opts); err != nil {
				s.sendToast(w, toast.WithErr(fmt.Errorf("failed to create namespace: %w", err)))
				return
			}
		}
		pkg = client.PackageBuilder(p.manifestName).
			WithVersion(p.version).
			WithRepositoryName(p.repositoryName).
			WithAutoUpdates(strings.ToLower(enableAutoUpdate) == "on").
			WithValues(values).
			WithNamespace(namespace).
			WithName(name).
			BuildPackage()
		err := install.NewInstaller(s.pkgClient).Install(ctx, pkg, opts)
		if err != nil {
			s.sendToast(w, toast.WithErr(fmt.Errorf("failed to install %v: %w", p.manifestName, err)))
		} else if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				s.sendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				s.sendYamlModal(w, yamlOutput, nil)
			}
		} else {
			s.swappingRedirect(w, "/packages", "main", "main")
			w.WriteHeader(http.StatusAccepted)
		}
	} else {
		pkg.Spec.PackageInfo.Version = p.version
		pkg.Spec.PackageInfo.RepositoryName = p.repositoryName
		pkg.Spec.Values = values
		opts := v1.UpdateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		if err := s.pkgClient.Packages(pkg.GetNamespace()).Update(ctx, pkg, opts); err != nil {
			s.sendToast(w, toast.WithErr(fmt.Errorf("failed to configure %v: %w", p.manifestName, err)))
			return
		}
		_, resolveErr := s.valueResolver.Resolve(ctx, values)
		if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				s.sendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				s.sendYamlModal(w, yamlOutput, resolveErr)
			}
		} else if resolveErr != nil {
			s.sendToast(w,
				toast.WithErr(fmt.Errorf("some values could not be resolved: %w", resolveErr)),
				toast.WithCssClass("warning"),
				toast.WithStatusCode(http.StatusAccepted))
		} else {
			s.sendToast(w, toast.WithMessage("Configuration updated successfully"))
		}
	}
}

// installOrConfigureClusterPackage either installs a new clusterpackage if it does not exist yet,
// or updates the configuration of an existing one.
// In either case, the parameters from the form are parsed and converted into ValueConfiguration objects, which are
// being set in the packages spec.
func (s *server) installOrConfigureClusterPackage(w http.ResponseWriter, r *http.Request, p *packageContextRequest) {
	ctx := r.Context()
	enableAutoUpdate := r.FormValue("enableAutoUpdate")
	dryRun, _ := strconv.ParseBool(r.FormValue("dryRun"))

	var err error
	pkg := &v1alpha1.ClusterPackage{}
	var mf *v1alpha1.PackageManifest
	if err = s.pkgClient.ClusterPackages().Get(ctx, p.manifestName, pkg); err != nil && !errors.IsNotFound(err) {
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to fetch clusterpackage %v: %w", p.manifestName, err)))
		return
	} else if err != nil {
		pkg = nil
	} else {
		// because disabled form elements are not submitted, we need to fall back on repo and version if advanced options disabled
		if p.repositoryName == "" {
			p.repositoryName = pkg.Spec.PackageInfo.RepositoryName
		}
		if p.version == "" {
			p.version = pkg.Spec.PackageInfo.Version
		}
	}

	mf, err = s.resolveManifest(ctx, pkg, p.repositoryName, p.manifestName, p.version)
	if repoerror.IsPartial(err) {
		fmt.Fprintf(os.Stderr, "problem fetching manifest and repo, but installation can continue: %v", err)
	} else if err != nil {
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to get manifest and repo of %v: %w", p.manifestName, err)))
		return
	}

	if values, err := extractValues(r, mf); err != nil {
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to parse values: %w", err)))
		return
	} else if pkg == nil {
		pkg = client.PackageBuilder(p.manifestName).
			WithVersion(p.version).
			WithRepositoryName(p.repositoryName).
			WithAutoUpdates(strings.ToLower(enableAutoUpdate) == "on").
			WithValues(values).
			BuildClusterPackage()
		opts := v1.CreateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		err := install.NewInstaller(s.pkgClient).Install(ctx, pkg, opts)
		if err != nil {
			s.sendToast(w, toast.WithErr(fmt.Errorf("failed to install %v: %w", p.manifestName, err)))
			return
		} else if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				s.sendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				s.sendYamlModal(w, yamlOutput, nil)
			}
		}
	} else {
		pkg.Spec.PackageInfo.Version = p.version
		pkg.Spec.PackageInfo.RepositoryName = p.repositoryName
		pkg.Spec.Values = values
		opts := v1.UpdateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		if err := s.pkgClient.ClusterPackages().Update(ctx, pkg, opts); err != nil {
			s.sendToast(w, toast.WithErr(fmt.Errorf("failed to configure %v: %w", p.manifestName, err)))
			return
		}
		_, resolveErr := s.valueResolver.Resolve(ctx, values)
		if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				s.sendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				s.sendYamlModal(w, yamlOutput, resolveErr)
			}
		} else if resolveErr != nil {
			s.sendToast(w,
				toast.WithErr(fmt.Errorf("some values could not be resolved: %w", resolveErr)),
				toast.WithCssClass("warning"),
				toast.WithStatusCode(http.StatusAccepted))
		} else {
			s.sendToast(w, toast.WithMessage("Configuration updated successfully"))
		}
	}
}

func (s *server) resolveManifest(ctx context.Context, pkg ctrlpkg.Package, repositoryName string, manifestName string, selectedVersion string) (
	*v1alpha1.PackageManifest, error) {

	var mf v1alpha1.PackageManifest
	var repoErr error
	if pkg.IsNil() ||
		(pkg.GetSpec().PackageInfo.RepositoryName != repositoryName || pkg.GetSpec().PackageInfo.Version != selectedVersion) {
		repoClient := s.repoClientset.ForRepoWithName(repositoryName)
		if err := repoClient.FetchPackageManifest(manifestName, selectedVersion, &mf); err != nil {
			return nil, multierr.Append(err, repoErr)
		}
	} else {
		if installedMf, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err != nil {
			return nil, multierr.Append(err, repoErr)
		} else {
			mf = *installedMf
		}
	}
	return &mf, repoErr
}
