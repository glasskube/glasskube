package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/glasskube/glasskube/api/v1alpha1"
	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/repo"
	repoerror "github.com/glasskube/glasskube/internal/repo/error"
	"github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/internal/web/components"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/cookie"
	opts "github.com/glasskube/glasskube/internal/web/options"
	"github.com/glasskube/glasskube/internal/web/responder"
	webtypes "github.com/glasskube/glasskube/internal/web/types"
	webutil "github.com/glasskube/glasskube/internal/web/util"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/describe"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/errors"
)

type packageDetailCommonData struct {
	Package              ctrlpkg.Package
	Status               *client.PackageStatus
	Manifest             *v1alpha1.PackageManifest
	PackageManifestUrl   string
	LatestVersion        string
	UpdateAvailable      bool
	ShowDiscussionLink   bool
	PackageHref          string
	AutoUpdaterInstalled bool
}

type packageDetailTemplateData struct {
	webtypes.TemplateContextHolder
	packageDetailCommonData
	ValidationResult     *dependency.ValidationResult
	ShowConflicts        bool
	SelectedVersion      string
	PackageIndex         *repo.PackageIndex
	Repositories         []v1alpha1.PackageRepository
	RepositoryName       string
	ShowConfiguration    bool
	ValueErrors          map[string]error
	DatalistOptions      map[string]*components.PkgConfigInputDatalistOptions
	AdvancedOptions      bool
	LostValueDefinitions []string
}

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

func getPackageContext(r *http.Request) *packageContext {
	return &packageContext{
		request: packageContextRequest{
			manifestName:   r.PathValue("manifestName"),
			namespace:      r.PathValue("namespace"),
			name:           r.PathValue("name"),
			repositoryName: r.FormValue("repositoryName"),
			version:        r.FormValue("version"),
			component:      r.FormValue("component"),
		},
	}
}

func GetPackageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := getPackageContext(r)

	var pkg *v1alpha1.Package
	var mf *v1alpha1.PackageManifest
	if p.request.namespaceAndNameSet() {
		var err error
		pkg, mf, err = describe.DescribeInstalledPackage(ctx, p.request.namespace, p.request.name)
		if err != nil && !errors.IsNotFound(err) {
			responder.SendToast(w, toast.WithErr(
				fmt.Errorf("failed to fetch installed package %v/%v: %w", p.request.namespace, p.request.name, err)))
			return
		} else if errors.IsNotFound(err) {
			responder.Redirect(w, "/packages")
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	p.pkg = pkg
	p.manifest = mf
	renderPackageDetailPage(w, r, p)
}

func GetClusterPackageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := getPackageContext(r)

	pkg, mf, err := describe.DescribeInstalledClusterPackage(ctx, p.request.manifestName)
	if err != nil && !errors.IsNotFound(err) {
		responder.SendToast(w,
			toast.WithErr(fmt.Errorf("failed to fetch installed clusterpackage %v: %w", p.request.manifestName, err)))
		return
	}
	p.pkg = pkg
	p.manifest = mf
	renderPackageDetailPage(w, r, p)
}

func renderPackageDetailPage(w http.ResponseWriter, r *http.Request, p *packageContext) {
	ctx := r.Context()

	pkgDetailCommonData, repos, idx, repoErr := resolvePkgDetailCommon(w, ctx, p)
	if pkgDetailCommonData == nil {
		// something happened, but already handled by resolvePkgDetailCommon func
		return
	}

	advancedOptions, err := cookie.GetAdvancedOptionsFromCookie(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get advanced options from cookie: %v\n", err)
	}

	if p.request.component == "header" {
		responder.SendComponent(w, r, "components/pkg-detail-header",
			responder.ContextualizedTemplate(&packageDetailTemplateData{
				packageDetailCommonData: *pkgDetailCommonData,
			}))
	} else {
		pkgClient := clicontext.PackageClientFromContext(ctx)
		repoClientset := clicontext.RepoClientsetFromContext(ctx)
		dependencyMgr := dependency.NewDependencyManager(
			clientadapter.NewPackageClientAdapter(pkgClient),
			repoClientset,
		)

		// note: not using updater here since it doesn't support changing repository or downgrading (and it probably shouldn't)
		// could maybe be refactored in the future (PackageChanger construct, which the "normal" updater is a special case of)
		validationResult := &dependency.ValidationResult{}
		var validationErr error
		var lostValueDefinitions []string
		valueErrors := make(map[string]error)
		datalistOptions := make(map[string]*components.PkgConfigInputDatalistOptions)

		if p.pkg.IsNil() {
			if p.manifest.Scope.IsCluster() {
				validationResult, validationErr =
					dependencyMgr.Validate(r.Context(), p.request.manifestName, "", p.manifest, p.request.version)
			} else {
				// In this case we don't know the actual namespace, but we can assume the default
				// TODO: make name and namespace depend on user input
				validationResult, validationErr =
					dependencyMgr.Validate(r.Context(), p.request.manifestName, p.manifest.DefaultNamespace, p.manifest, p.request.version)
			}
		} else if shouldMigrateManifest(p) {
			validationResult, validationErr =
				dependencyMgr.Validate(r.Context(), p.pkg.GetName(), p.pkg.GetNamespace(), p.manifest, p.request.version)
		}
		if validationErr != nil {
			responder.SendToast(w,
				toast.WithErr(fmt.Errorf("failed to validate dependencies of %v (%v): %w",
					p.request.manifestName, p.request.version, validationErr)))
			return
		}

		nsOptions, _ := opts.GetNamespaceOptions(ctx)
		if !p.pkg.IsNil() {
			pkgsOptions, _ := opts.GetPackagesOptions(r.Context())
			for key, v := range p.pkg.GetSpec().Values {
				k8sClient := clicontext.KubernetesClientFromContext(ctx)
				valueResolver := manifestvalues.NewResolver(
					clientadapter.NewPackageClientAdapter(pkgClient),
					clientadapter.NewKubernetesClientAdapter(k8sClient),
				)
				if resolved, err := valueResolver.ResolveValue(r.Context(), v); err != nil {
					valueErrors[key] = util.GetRootCause(err)
				} else if valDef, exists := p.manifest.ValueDefinitions[key]; !exists {
					// can happen when a different repo/version of the pkg is requested (advanced options)
					lostValueDefinitions = append(lostValueDefinitions, key)
				} else if err := manifestvalues.ValidateSingle(key, valDef, resolved); err != nil {
					valueErrors[key] = err
				}
				if v.ValueFrom != nil {
					options, err := opts.GetDatalistOptions(ctx, v.ValueFrom, nsOptions, pkgsOptions)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%v\n", err)
					}
					datalistOptions[key] = options
				}
			}
		}
		datalistOptions[""] = &components.PkgConfigInputDatalistOptions{Namespaces: nsOptions}

		templateData := &packageDetailTemplateData{
			packageDetailCommonData: *pkgDetailCommonData,
			ValidationResult:        validationResult,
			ShowConflicts:           validationResult.Status == dependency.ValidationResultStatusConflict,
			SelectedVersion:         p.request.version,
			PackageIndex:            &idx,
			Repositories:            repos,
			RepositoryName:          p.request.repositoryName,
			ShowConfiguration:       (!p.pkg.IsNil() && len(p.manifest.ValueDefinitions) > 0 && p.pkg.GetDeletionTimestamp().IsZero()) || p.pkg.IsNil(),
			ValueErrors:             valueErrors,
			DatalistOptions:         datalistOptions,
			AdvancedOptions:         advancedOptions,
			LostValueDefinitions:    lostValueDefinitions,
		}
		responder.SendPage(w, r, "pages/package", responder.ContextualizedTemplate(templateData), responder.WithPartialErr(repoErr))
	}
}

func resolvePkgDetailCommon(w http.ResponseWriter, ctx context.Context, p *packageContext) (*packageDetailCommonData, []v1alpha1.PackageRepository, repo.PackageIndex, error) {
	if !p.pkg.IsNil() {
		// for installed packages, the installed repo + version is the fallback, if they are not requested explicitly
		if p.request.repositoryName == "" {
			p.request.repositoryName = p.pkg.GetSpec().PackageInfo.RepositoryName
		}
		if p.request.version == "" {
			p.request.version = p.pkg.GetSpec().PackageInfo.Version
		}
	}

	var repoErr error
	var repos []v1alpha1.PackageRepository
	var usedRepo *v1alpha1.PackageRepository
	if p.request.repositoryName, repos, usedRepo, repoErr = resolveRepos(
		ctx, p.request.manifestName, p.request.repositoryName); repoerror.IsComplete(repoErr) {
		responder.SendToast(w, toast.WithErr(repoErr))
		return nil, nil, repo.PackageIndex{}, nil
	}

	var idx repo.PackageIndex
	var latestVersion string
	var err error
	if idx, latestVersion, p.request.version, err = resolveVersions(
		ctx, p.request.repositoryName, p.request.manifestName, p.request.version); err != nil {
		responder.SendToast(w,
			toast.WithErr(fmt.Errorf("failed to fetch package index of %v in repo %v: %w",
				p.request.manifestName, p.request.repositoryName, multierr.Append(repoErr, err))))
		return nil, nil, repo.PackageIndex{}, nil
	}

	packageManifestUrl := ""
	repoClientset := clicontext.RepoClientsetFromContext(ctx)
	repoClient := repoClientset.ForRepo(*usedRepo)
	if p.manifest == nil || shouldMigrateManifest(p) {
		p.manifest = &v1alpha1.PackageManifest{}
		if err := repoClient.FetchPackageManifest(p.request.manifestName, p.request.version, p.manifest); err != nil {
			responder.SendToast(w,
				toast.WithErr(fmt.Errorf("failed to fetch manifest of %v (%v) in repo %v: %w",
					p.request.manifestName, p.request.version, p.request.repositoryName, err)))
			return nil, nil, repo.PackageIndex{}, nil
		}
	}
	if url, err := repoClient.GetPackageManifestURL(p.request.manifestName, p.request.version); err != nil {
		fmt.Fprintf(os.Stderr, "failed to get package manifest url of %v (%v) in repo %v: %v",
			p.request.manifestName, p.request.version, p.request.repositoryName, err)
	} else {
		packageManifestUrl = url
	}

	autoUpdaterInstalled, err := clientutils.IsAutoUpdaterInstalled(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to check whether auto updater is installed: %v\n", err)
	}

	return &packageDetailCommonData{
		Package:              p.pkg,
		Status:               client.GetStatusOrPending(p.pkg),
		Manifest:             p.manifest,
		PackageManifestUrl:   packageManifestUrl,
		LatestVersion:        latestVersion,
		UpdateAvailable:      isUpdateAvailableForPkg(ctx, p.pkg),
		ShowDiscussionLink:   usedRepo.IsGlasskubeRepo(),
		PackageHref:          webutil.GetPackageHrefWithFallback(p.pkg, p.manifest),
		AutoUpdaterInstalled: autoUpdaterInstalled,
	}, repos, idx, repoErr
}

func shouldMigrateManifest(p *packageContext) bool {
	if !p.pkg.IsNil() {
		return p.request.repositoryName != p.pkg.GetSpec().PackageInfo.RepositoryName ||
			p.request.version != p.pkg.GetSpec().PackageInfo.Version
	}
	return false
}

func resolveVersions(ctx context.Context, repositoryName string, pkgName string, selectedVersion string) (repo.PackageIndex, string, string, error) {
	var idx repo.PackageIndex
	repoClientset := clicontext.RepoClientsetFromContext(ctx)
	if err := repoClientset.ForRepoWithName(repositoryName).FetchPackageIndex(pkgName, &idx); err != nil {
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

func resolveRepos(ctx context.Context, manifestName string, repositoryName string) (
	string, []v1alpha1.PackageRepository, *v1alpha1.PackageRepository, error) {
	var repos []v1alpha1.PackageRepository
	var err error
	repoClientset := clicontext.RepoClientsetFromContext(ctx)
	pkgClient := clicontext.PackageClientFromContext(ctx)

	if repos, err = repoClientset.Meta().GetReposForPackage(manifestName); err != nil {
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
	if err := pkgClient.PackageRepositories().Get(ctx, repositoryName, &usedRepo); err != nil {
		return "", nil, nil, err
	}

	return repositoryName, repos, &usedRepo, err
}
