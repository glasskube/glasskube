package web

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/repo"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/internal/web/components/pkg_config_input"
	"github.com/glasskube/glasskube/internal/web/handler"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/manifest"
	"github.com/glasskube/glasskube/pkg/open"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/uninstall"
	"github.com/glasskube/glasskube/pkg/update"
	"github.com/gorilla/mux"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
)

//go:embed root
//go:embed templates
var embeddedFs embed.FS
var webFs fs.FS = embeddedFs

func init() {
	if config.IsDevBuild() {
		if _, err := os.Lstat(templatesBaseDir); err == nil {
			webFs = os.DirFS(templatesBaseDir)
		}
	}
}

type ServerOptions struct {
	Host               string
	Port               int32
	Kubeconfig         string
	LogLevel           int
	SkipOpeningBrowser bool
}

func NewServer(options ServerOptions) *server {
	server := server{
		ServerOptions:      options,
		configLoader:       &defaultConfigLoader{options.Kubeconfig},
		forwarders:         make(map[string]*open.OpenResult),
		updateTransactions: make(map[int]update.UpdateTransaction),
		templates:          templates{},
	}
	return &server
}

type server struct {
	ServerOptions
	configLoader
	listener           net.Listener
	restConfig         *rest.Config
	rawConfig          *api.Config
	pkgClient          client.PackageV1Alpha1Client
	repoClientset      repoclient.RepoClientset
	k8sClient          *kubernetes.Clientset
	sseHub             *SSEHub
	namespaceLister    *corev1.NamespaceLister
	configMapLister    *corev1.ConfigMapLister
	secretLister       *corev1.SecretLister
	forwarders         map[string]*open.OpenResult
	dependencyMgr      *dependency.DependendcyManager
	updateMutex        sync.Mutex
	updateTransactions map[int]update.UpdateTransaction
	valueResolver      *manifestvalues.Resolver
	isBootstrapped     bool
	templates          templates
}

func (s *server) RestConfig() *rest.Config {
	return s.restConfig
}

func (s *server) RawConfig() *api.Config {
	return s.rawConfig
}

func (s *server) Client() client.PackageV1Alpha1Client {
	return s.pkgClient
}

func (s *server) K8sClient() *kubernetes.Clientset {
	return s.k8sClient
}

func (s *server) RepoClient() repoclient.RepoClientset {
	return s.repoClientset
}

func initLogging(level int) {
	klog.InitFlags(nil)
	_ = flag.Set("v", strconv.Itoa(level))
	flag.Parse()
}

func (s *server) Start(ctx context.Context) error {
	if s.listener != nil {
		return errors.New("server is already listening")
	}

	if s.LogLevel != 0 {
		initLogging(s.LogLevel)
	} else if config.IsDevBuild() {
		initLogging(5)
	}

	s.templates.parseTemplates()
	if config.IsDevBuild() {
		if err := s.templates.watchTemplates(); err != nil {
			fmt.Fprintf(os.Stderr, "templates will not be parsed after changes: %v\n", err)
		}
	}
	s.sseHub = NewHub()
	_ = s.ensureBootstrapped(ctx)

	root, err := fs.Sub(webFs, "root")
	if err != nil {
		return err
	}

	fileServer := http.FileServer(http.FS(root))

	router := mux.NewRouter()
	router.Use(telemetry.HttpMiddleware())
	router.PathPrefix("/static/").Handler(fileServer)
	router.Handle("/favicon.ico", fileServer)
	router.HandleFunc("/events", s.sseHub.handler)
	router.HandleFunc("/support", s.supportPage)
	router.HandleFunc("/kubeconfig", s.kubeconfigPage)
	router.Handle("/bootstrap", s.requireKubeconfig(s.bootstrapPage))
	router.Handle("/kubeconfig/persist", s.requireKubeconfig(s.persistKubeconfig))
	router.Handle("/packages", s.requireReady(s.packages))
	router.Handle("/packages/update", s.requireReady(s.update))
	router.Handle("/packages/update/modal", s.requireReady(s.updateModal))
	router.Handle("/packages/uninstall", s.requireReady(s.uninstall))
	router.Handle("/packages/uninstall/modal", s.requireReady(s.uninstallModal))
	router.Handle("/packages/open", s.requireReady(s.open))
	router.Handle("/packages/{pkgName}", s.requireReady(s.packageDetail))
	router.Handle("/packages/{pkgName}/discussion", s.requireReady(s.packageDiscussion))
	router.Handle("/packages/{pkgName}/discussion/badge", s.requireReady(s.discussionBadge))
	router.Handle("/packages/{pkgName}/configure", s.requireReady(s.installOrConfigurePackage))
	router.Handle("/packages/{pkgName}/configure/advanced", s.requireReady(s.advancedConfiguration))
	router.Handle("/packages/{pkgName}/configuration/{valueName}", s.requireReady(s.packageConfigurationInput))
	router.Handle("/packages/{pkgName}/configuration/{valueName}/datalists/names", s.requireReady(s.namesDatalist))
	router.Handle("/packages/{pkgName}/configuration/{valueName}/datalists/keys", s.requireReady(s.keysDatalist))
	router.Handle("/settings", s.requireReady(s.settingsPage))
	router.Handle("/settings/repository/{repoName}", s.requireReady(s.repositoryconfig))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/packages", http.StatusFound)
	})
	http.Handle("/", s.enrichContext(router))

	bindAddr := fmt.Sprintf("%v:%d", s.Host, s.Port)

	s.listener, err = net.Listen("tcp", bindAddr)
	if err != nil {
		// Checks if Port Conflict Error exists
		if isPortConflictError(err) {
			userInput := cliutils.YesNoPrompt(
				"Port is already in use.\nShould glasskube use a different port? (Y/n): ", true)
			if userInput {
				s.listener, err = net.Listen("tcp", ":0")
				if err != nil {
					panic(err)
				}
				bindAddr = fmt.Sprintf("%v:%d", s.Host, s.listener.Addr().(*net.TCPAddr).Port)
			} else {
				fmt.Println("Exiting. User chose not to use a different port.")
				cliutils.ExitWithError()
			}
		} else {
			// If no Port Conflict error is found, return other errors
			return err
		}
	}

	fmt.Printf("glasskube UI is available at http://%v\n", bindAddr)
	if !s.SkipOpeningBrowser {
		_ = cliutils.OpenInBrowser("http://" + bindAddr)
	}

	go s.sseHub.Run()
	server := &http.Server{}
	err = server.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *server) updateModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := r.FormValue("packageName")
	pkgs := make([]string, 0, 1)
	if pkgName != "" {
		pkgs = append(pkgs, pkgName)
	}

	updates := make([]map[string]any, 0)
	updater := update.NewUpdater(ctx).WithStatusWriter(statuswriter.Stderr())
	ut, err := updater.Prepare(ctx, pkgs)
	if err != nil {
		s.respondAlertAndLog(w, err, "An error occurred preparing update of "+pkgName, "danger")
		return
	}
	utId := rand.Int()
	s.updateMutex.Lock()
	s.updateTransactions[utId] = *ut
	s.updateMutex.Unlock()

	for _, u := range ut.Items {
		if u.UpdateRequired() {
			updates = append(updates, map[string]any{
				"Name":           u.Package.Name,
				"CurrentVersion": u.Package.Spec.PackageInfo.Version,
				"LatestVersion":  u.Version,
			})
		}
	}
	for _, req := range ut.Requirements {
		updates = append(updates, map[string]any{
			"Name":           req.Name,
			"CurrentVersion": "-",
			"LatestVersion":  req.Version,
		})
	}

	err = s.templates.pkgUpdateModalTmpl.Execute(w, map[string]any{
		"UpdateTransactionId": utId,
		"Updates":             updates,
		"PackageName":         pkgName,
	})
	checkTmplError(err, "pkgUpdateModalTmpl")
}

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	updater := update.NewUpdater(ctx).WithStatusWriter(statuswriter.Stderr())
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()
	utIdStr := r.FormValue("updateTransactionId")
	if utId, err := strconv.Atoi(utIdStr); err != nil {
		s.respondAlertAndLog(w, err, "Failed to parse updateTransactionId", "danger")
		return
	} else if ut, ok := s.updateTransactions[utId]; !ok {
		s.respondAlert(w, fmt.Sprintf("Failed to find UpdateTransaction with ID %d", utId), "danger")
		return
	} else if _, err := updater.Apply(ctx, &ut); err != nil {
		delete(s.updateTransactions, utId)
		s.respondAlertAndLog(w, err, "An error occurred during the update", "danger")
		return
	} else {
		delete(s.updateTransactions, utId)
	}
}

func (s *server) uninstallModal(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	var pruned []string
	var err error
	if g, err1 := s.dependencyMgr.NewGraph(r.Context()); err1 != nil {
		err = fmt.Errorf("error validating uninstall: %w", err1)
	} else {
		g.Delete(pkgName)
		pruned = g.Prune()
		if err1 := g.Validate(); err1 != nil {
			err = fmt.Errorf("%v cannot be uninstalled: %w", pkgName, err1)
		}
	}
	err = s.templates.pkgUninstallModalTmpl.Execute(w, map[string]any{
		"PackageName": pkgName,
		"Pruned":      pruned,
		"Err":         err,
	})
	checkTmplError(err, "pkgUninstallModalTmpl")
}

func (s *server) uninstall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := r.FormValue("packageName")
	var pkg v1alpha1.ClusterPackage
	if err := s.pkgClient.ClusterPackages().Get(ctx, pkgName, &pkg); err != nil {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching %v during uninstall", pkgName), "danger")
		return
	}
	if err := uninstall.NewUninstaller(s.pkgClient).
		WithStatusWriter(statuswriter.Stderr()).
		Uninstall(ctx, &pkg); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred uninstalling "+pkgName, "danger")
		return
	}
}

func (s *server) open(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	if result, ok := s.forwarders[pkgName]; ok {
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		return
	}

	var pkg v1alpha1.ClusterPackage
	if err := s.pkgClient.ClusterPackages().Get(r.Context(), pkgName, &pkg); err != nil {
		s.respondAlertAndLog(w, err, "Could not get ClusterPackage", "danger")
		return
	}

	result, err := open.NewOpener().Open(r.Context(), &pkg, "", 0)
	if err != nil {
		s.respondAlertAndLog(w, err, "Could not open "+pkgName, "danger")
	} else {
		s.forwarders[pkgName] = result
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *server) packages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	packages, listErr := list.NewLister(ctx).GetPackagesWithStatus(ctx, list.ListOptions{IncludePackageInfos: true})
	if listErr != nil && len(packages) == 0 {
		listErr = fmt.Errorf("could not load packages: %w", listErr)
		fmt.Fprintf(os.Stderr, "%v\n", listErr)
	}

	// Call isUpdateAvailable for each installed package.
	// This is not the same as getting all updates in a single transaction, because some dependency
	// conflicts could be resolvable by installing individual packages.
	packageUpdateAvailable := map[string]bool{}
	for _, pkg := range packages {
		packageUpdateAvailable[pkg.Name] = pkg.Package != nil && s.isUpdateAvailable(r.Context(), pkg.Name)
	}

	tmplErr := s.templates.pkgsPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Packages":               packages,
		"PackageUpdateAvailable": packageUpdateAvailable,
		"UpdatesAvailable":       s.isUpdateAvailable(r.Context()),
	}, listErr))
	checkTmplError(tmplErr, "packages")
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")

	pkg, manifest, err := describe.DescribeInstalledPackage(ctx, pkgName)
	if err != nil && !apierrors.IsNotFound(err) {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching package details of installed package %v", pkgName),
			"danger")
		return
	} else if pkg != nil {
		repositoryName = pkg.Spec.PackageInfo.RepositoryName
	}

	var repos []v1alpha1.PackageRepository
	if repos, err = s.repoClientset.Meta().GetReposForPackage(pkgName); err != nil {
		fmt.Fprintf(os.Stderr, "error getting repos for package; %v", err)
	} else if repositoryName == "" {
		if len(repos) == 0 {
			s.respondAlertAndLog(w, fmt.Errorf("%v not found in any repository", pkgName), "", "danger")
			return
		}
		for _, r := range repos {
			repositoryName = r.Name
			if r.IsDefaultRepository() {
				break
			}
		}
	}

	var usedRepo v1alpha1.PackageRepository
	if err := s.pkgClient.PackageRepositories().Get(r.Context(), repositoryName, &usedRepo); err != nil {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching repository %v", repositoryName),
			"danger")
		return
	}

	var idx repo.PackageIndex
	if err := s.repoClientset.ForRepoWithName(repositoryName).FetchPackageIndex(pkgName, &idx); err != nil {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching package index of %v in repository %v", pkgName, repositoryName),
			"danger")
		return
	}
	latestVersion := idx.LatestVersion

	if selectedVersion == "" {
		selectedVersion = latestVersion
	} else if !slices.ContainsFunc(idx.Versions, func(item types.PackageIndexItem) bool {
		return item.Version == selectedVersion
	}) {
		selectedVersion = latestVersion
	}

	if manifest == nil {
		manifest = &v1alpha1.PackageManifest{}
		if err := s.repoClientset.ForRepoWithName(repositoryName).
			FetchPackageManifest(pkgName, selectedVersion, manifest); err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching manifest of %v in version %v in repository %v",
					pkgName, selectedVersion, repositoryName),
				"danger")
			return
		}
	}

	res, err := s.dependencyMgr.Validate(r.Context(), manifest, selectedVersion)
	if err != nil {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred validating dependencies of %v in version %v", pkgName, selectedVersion),
			"danger")
		return
	}

	valueErrors := make(map[string]error)
	datalistOptions := make(map[string]*pkg_config_input.PkgConfigInputDatalistOptions)
	if pkg != nil {
		nsOptions, _ := s.getNamespaceOptions()
		pkgsOptions, _ := s.getPackagesOptions(r.Context())
		for key, v := range pkg.Spec.Values {
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
		"Package":            pkg,
		"Status":             client.GetStatusOrPending(pkg),
		"Manifest":           manifest,
		"LatestVersion":      latestVersion,
		"UpdateAvailable":    pkg != nil && s.isUpdateAvailable(r.Context(), pkgName),
		"AutoUpdate":         clientutils.AutoUpdateString(pkg, "Disabled"),
		"ValidationResult":   res,
		"ShowConflicts":      res.Status == dependency.ValidationResultStatusConflict,
		"SelectedVersion":    selectedVersion,
		"PackageIndex":       &idx,
		"Repositories":       repos,
		"RepositoryName":     repositoryName,
		"ShowConfiguration":  (pkg != nil && len(manifest.ValueDefinitions) > 0 && pkg.DeletionTimestamp.IsZero()) || pkg == nil,
		"ValueErrors":        valueErrors,
		"DatalistOptions":    datalistOptions,
		"ShowDiscussionLink": usedRepo.IsGlasskubeRepo(),
	}, err))
	checkTmplError(err, fmt.Sprintf("package-detail (%s)", pkgName))
}

// installOrConfigurePackage is an endpoint which takes POST requests, containing all necessary parameters to either
// install a new package if it does not exist yet, or update the configuration of an existing package.
// The name of the concerned package is given in the pkgName query parameter.
// In case the given package is not installed yet in the cluster, there must be a form parameter selectedVersion
// containing which version should be installed.
// In either case, the parameters from the form are parsed and converted into ValueConfiguration objects, which are
// being set in the packages spec.
func (s *server) installOrConfigurePackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	enableAutoUpdate := r.FormValue("enableAutoUpdate")
	pkg := &v1alpha1.ClusterPackage{}
	var mf v1alpha1.PackageManifest
	if err := s.pkgClient.ClusterPackages().Get(ctx, pkgName, pkg); err != nil && !apierrors.IsNotFound(err) {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching package details of %v", pkgName), "danger")
		return
	} else if err != nil {
		pkg = nil
	}

	if pkg == nil {
		var repoClient repoclient.RepoClient
		if len(repositoryName) == 0 {
			repos, err := s.repoClientset.Meta().GetReposForPackage(pkgName)
			if err != nil {
				s.respondAlertAndLog(w, err, "", "danger")
				return
			}
			switch len(repos) {
			case 0:
				// TODO: show error in UI
				fmt.Fprintf(os.Stderr, "package not found in any repository")
				return
			case 1:
				repositoryName = repos[0].Name
				repoClient = s.repoClientset.ForRepo(repos[0])
			default:
				// TODO: show error in UI
				fmt.Fprintf(os.Stderr, "package found in multiple repositories")
				return
			}
		} else {
			repoClient = s.repoClientset.ForRepoWithName(repositoryName)
		}
		if err := repoClient.FetchPackageManifest(pkgName, selectedVersion, &mf); err != nil {
			s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching manifest of %v in version %v", pkgName, selectedVersion), "danger")
			return
		}
	} else {
		if mf1, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err != nil {
			s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching package details of %v", pkgName), "danger")
			return
		} else {
			mf = *mf1
		}
	}

	if values, err := extractValues(r, &mf); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred parsing the form", "danger")
		return
	} else if pkg == nil {
		pkg = client.ClusterPackageBuilder(pkgName).
			WithVersion(selectedVersion).
			WithRepositoryName(repositoryName).
			WithAutoUpdates(strings.ToLower(enableAutoUpdate) == "on").
			WithValues(values).
			Build()
		opts := metav1.CreateOptions{}
		err := install.NewInstaller(s.pkgClient).
			WithStatusWriter(statuswriter.Stderr()).
			Install(ctx, pkg, opts)
		if err != nil {
			s.respondAlertAndLog(w, err, "An error occurred installing "+pkgName, "danger")
			return
		}
	} else {
		pkg.Spec.Values = values
		if err := s.pkgClient.ClusterPackages().Update(ctx, pkg); err != nil {
			s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred updating package %v", pkgName), "danger")
			return
		}
		if _, err := s.valueResolver.Resolve(ctx, values); err != nil {
			s.respondAlertAndLog(w, err, "Some values could not be resolved: ", "warning")
		} else {
			err := s.templates.alertTmpl.Execute(w, map[string]any{
				"Message":     "Configuration updated successfully",
				"Dismissible": true,
				"Type":        "success",
			})
			checkTmplError(err, "success")
		}
	}
}

// advancedConfiguration is a GET+POST endpoint which can be used for advanced package installation options, most notably
// for changing the package repository and changing to a specific (maybe even lower than installed) version of the package.
// It is only intended to be used for already installed packages, for new packages these options exist anyway and
// should be available for every user.
func (s *server) advancedConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	pkg, manifest, err := describe.DescribeInstalledPackage(ctx, pkgName)
	if err != nil && !apierrors.IsNotFound(err) {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching package details of installed package %v", pkgName),
			"danger")
		return
	} else if pkg == nil {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("Package %v is not installed", pkgName),
			"danger")
		return
	} else if repositoryName == "" {
		repositoryName = pkg.Spec.PackageInfo.RepositoryName
	}
	var repos []v1alpha1.PackageRepository
	if repos, err = s.repoClientset.Meta().GetReposForPackage(pkgName); err != nil {
		fmt.Fprintf(os.Stderr, "error getting repos for package; %v", err)
	} else if repositoryName == "" {
		if len(repos) == 0 {
			s.respondAlertAndLog(w, fmt.Errorf("%v not found in any repository", pkgName), "", "danger")
			return
		}
		for _, r := range repos {
			repositoryName = r.Name
			if r.IsDefaultRepository() {
				break
			}
		}
	}

	if r.Method == http.MethodGet {
		var idx repo.PackageIndex
		if err := s.repoClientset.ForRepoWithName(repositoryName).FetchPackageIndex(pkgName, &idx); err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching package index of %v in repository %v", pkgName, repositoryName),
				"danger")
			return
		}
		latestVersion := idx.LatestVersion

		if selectedVersion == "" {
			selectedVersion = latestVersion
		} else if !slices.ContainsFunc(idx.Versions, func(item types.PackageIndexItem) bool {
			return item.Version == selectedVersion
		}) {
			selectedVersion = latestVersion
		}

		res, err := s.dependencyMgr.Validate(r.Context(), manifest, selectedVersion)
		if err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred validating dependencies of %v in version %v", pkgName, selectedVersion),
				"danger")
			return
		}

		err = s.templates.pkgConfigAdvancedTmpl.Execute(w, s.enrichPage(r, map[string]any{
			"Status":           client.GetStatusOrPending(pkg),
			"Manifest":         manifest,
			"LatestVersion":    latestVersion,
			"ValidationResult": res,
			"ShowConflicts":    res.Status == dependency.ValidationResultStatusConflict,
			"SelectedVersion":  selectedVersion,
			"PackageIndex":     &idx,
			"Repositories":     repos,
			"RepositoryName":   repositoryName,
		}, err))
		checkTmplError(err, fmt.Sprintf("advanced-config (%s)", pkgName))
	} else if r.Method == http.MethodPost {
		pkg.Spec.PackageInfo.Version = selectedVersion
		if repositoryName != "" {
			pkg.Spec.PackageInfo.RepositoryName = repositoryName
		}
		if err := s.pkgClient.ClusterPackages().Update(ctx, pkg); err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred updating package %v to version %v in repo %v", pkgName, selectedVersion, repositoryName),
				"danger")
			return
		} else {
			err := s.templates.alertTmpl.Execute(w, map[string]any{
				"Message":     "Configuration updated successfully",
				"Dismissible": true,
				"Type":        "success",
			})
			checkTmplError(err, "success")
		}
	}
}

func (s *server) supportPage(w http.ResponseWriter, r *http.Request) {
	if err := s.ensureBootstrapped(r.Context()); err != nil {
		if err.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
			return
		}
		err := s.templates.supportPageTmpl.Execute(w, &map[string]any{
			"CurrentContext":            "",
			"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
			"Err":                       err,
		})
		checkTmplError(err, "support")
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *server) bootstrapPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == "POST" {
		client := bootstrap.NewBootstrapClient(s.restConfig)
		if _, err := client.Bootstrap(ctx, bootstrap.DefaultOptions()); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			err := s.templates.bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-failure", nil)
			checkTmplError(err, "bootstrap-failure")
		} else {
			err := s.templates.bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-success", nil)
			checkTmplError(err, "bootstrap-success")
		}
	} else {
		isBootstrapped, err := bootstrap.IsBootstrapped(ctx, s.restConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
		} else if isBootstrapped {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		tplErr := s.templates.bootstrapPageTmpl.Execute(w, &map[string]any{
			"CloudId":        telemetry.GetMachineId(),
			"CurrentContext": s.rawConfig.CurrentContext,
			"Err":            err,
		})
		checkTmplError(tplErr, "bootstrap")
	}
}

func (s *server) kubeconfigPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("kubeconfig")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.loadBytesConfig(data)
		if err := s.checkKubeconfig(); err != nil {
			fmt.Fprintf(os.Stderr, "The selected kubeconfig is invalid: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "The selected kubeconfig is valid!")
		}
	}

	configErr := s.checkKubeconfig()
	var currentContext string
	if s.rawConfig != nil {
		currentContext = s.rawConfig.CurrentContext
	}
	tplErr := s.templates.kubeconfigPageTmpl.Execute(w, map[string]any{
		"CloudId":                   telemetry.GetMachineId(),
		"CurrentContext":            currentContext,
		"ConfigErr":                 configErr,
		"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
		"DefaultKubeconfigExists":   defaultKubeconfigExists(),
	})
	checkTmplError(tplErr, "kubeconfig")
}

func (s *server) settingsPage(w http.ResponseWriter, r *http.Request) {
	var repos v1alpha1.PackageRepositoryList
	if err := s.pkgClient.PackageRepositories().GetAll(r.Context(), &repos); err != nil {
		s.respondAlertAndLog(w, err, "Failed to fetch repositories", "danger")
		return
	}

	tmplErr := s.templates.settingsPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Repositories": repos.Items,
	}, nil))
	checkTmplError(tmplErr, "settings")
}

func (s *server) enrichPage(r *http.Request, data map[string]any, err error) map[string]any {
	data["CloudId"] = telemetry.GetMachineId()
	if pathParts := strings.Split(r.URL.Path, "/"); len(pathParts) >= 2 {
		data["NavbarActiveItem"] = pathParts[1]
	}
	data["Error"] = err
	data["CurrentContext"] = s.rawConfig.CurrentContext
	operatorVersion, clientVersion, err := s.getGlasskubeVersions(r.Context())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check for version mismatch: %v\n", err)
	} else if operatorVersion != nil && clientVersion != nil && !operatorVersion.Equal(clientVersion) {
		data["VersionMismatchWarning"] = true
	}
	if operatorVersion != nil && clientVersion != nil && !config.IsDevBuild() {
		data["VersionDetails"] = map[string]any{
			"OperatorVersion":     operatorVersion.String(),
			"ClientVersion":       clientVersion.String(),
			"NeedsOperatorUpdate": operatorVersion.LessThan(clientVersion),
		}
	}
	if config.IsDevBuild() {
		data["VersionDetails"] = map[string]any{
			"OperatorVersion": config.Version,
			"ClientVersion":   config.Version,
		}
	}
	data["CacheBustingString"] = config.Version
	return data
}

func (server *server) getGlasskubeVersions(ctx context.Context) (*semver.Version, *semver.Version, error) {
	if !config.IsDevBuild() {
		if operatorVersion, err := clientutils.GetPackageOperatorVersion(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check package operator version: %v\n", err)
			return nil, nil, err
		} else if parsedOperator, err := semver.NewVersion(operatorVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse operator version %v: %v\n", operatorVersion, err)
			return nil, nil, err
		} else if parsedClient, err := semver.NewVersion(config.Version); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse client version %v: %v\n", config.Version, err)
			return nil, nil, err
		} else {
			return parsedOperator, parsedClient, nil
		}
	}
	return nil, nil, nil
}

func (s *server) persistKubeconfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !defaultKubeconfigExists() {
			if err := clientcmd.WriteToFile(*s.rawConfig, clientcmd.RecommendedHomeFile); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Redirect(w, r, "/", http.StatusFound)
			}
		} else {
			fmt.Fprintln(os.Stderr, "default kubeconfig already exists! nothing was saved")
			http.Error(w, "", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "only POST is supported", http.StatusMethodNotAllowed)
	}
}

func (server *server) loadBytesConfig(data []byte) {
	server.configLoader = &bytesConfigLoader{data}
}

func (server *server) checkKubeconfig() ServerConfigError {
	if server.pkgClient == nil {
		return server.initKubeConfig()
	} else {
		return nil
	}
}

// ensureBootstrapped checks for a valid kubeconfig (see checkKubeconfig), and whether glasskube is bootstrapped in
// the given cluster. If either of these checks fail, a ServerConfigError is returned, otherwise the result of the
// check is cached in isBootstrapped and the check will not run anymore after that. After the first successful check,
// additional components are intialized (which can only be done once glasskube is known to be bootstrapped) –
// see initWhenBootstrapped
func (server *server) ensureBootstrapped(ctx context.Context) ServerConfigError {
	if server.isBootstrapped {
		return nil
	}
	if err := server.checkKubeconfig(); err != nil {
		return err
	}

	isBootstrapped, err := bootstrap.IsBootstrapped(ctx, server.restConfig)
	if !isBootstrapped || err != nil {
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
		}
		return newBootstrapErr(err)
	}
	server.isBootstrapped = isBootstrapped
	server.initWhenBootstrapped(ctx)
	return nil
}

func (server *server) initKubeConfig() ServerConfigError {
	restConfig, rawConfig, err := server.LoadConfig()
	if err != nil {
		return newKubeconfigErr(err)
	}
	client, err := client.New(restConfig)
	if err != nil {
		return newKubeconfigErr(err)
	}
	telemetry.InitClient(restConfig)

	server.restConfig = restConfig
	server.rawConfig = rawConfig
	server.pkgClient = client // be aware that server.pkgClient is overridden with the cached client once bootstrap check succeeded
	return nil
}

func (server *server) initWhenBootstrapped(ctx context.Context) {
	server.k8sClient = kubernetes.NewForConfigOrDie(server.restConfig)
	server.initCachedClient(context.WithoutCancel(ctx))
	server.repoClientset = repoclient.NewClientset(
		clientadapter.NewPackageClientAdapter(server.pkgClient),
		clientadapter.NewKubernetesClientAdapter(server.k8sClient),
	)
	server.templates.repoClientset = server.repoClientset
	server.dependencyMgr = dependency.NewDependencyManager(
		clientadapter.NewPackageClientAdapter(server.pkgClient),
		server.repoClientset,
	)
	server.valueResolver = manifestvalues.NewResolver(
		clientadapter.NewPackageClientAdapter(server.pkgClient),
		clientadapter.NewKubernetesClientAdapter(server.k8sClient),
	)
	factory := informers.NewSharedInformerFactory(server.k8sClient, 0)
	c := make(chan struct{})
	namespaceLister := factory.Core().V1().Namespaces().Lister()
	server.namespaceLister = &namespaceLister
	configMapLister := factory.Core().V1().ConfigMaps().Lister()
	server.configMapLister = &configMapLister
	secretLister := factory.Core().V1().Secrets().Lister()
	server.secretLister = &secretLister
	factory.Start(c)
}

func (server *server) initCachedClient(ctx context.Context) {
	packageStore, packageController := server.initPackageStoreAndController(ctx)
	packageInfoStore, packageInfoController := server.initPackageInfoStoreAndController(ctx)
	packageRepoStore, packageRepoController := server.initPackageRepoStoreAndController(ctx)
	go packageController.Run(ctx.Done())
	go packageInfoController.Run(ctx.Done())
	go packageRepoController.Run(ctx.Done())
	server.pkgClient = server.pkgClient.WithStores(packageStore, packageInfoStore, packageRepoStore)
}

func (s *server) enrichContext(h http.Handler) http.Handler {
	return &handler.ContextEnrichingHandler{Source: s, Handler: h}
}

func (s *server) requireReady(h http.HandlerFunc) http.Handler {
	return &handler.PreconditionHandler{
		Precondition: func(r *http.Request) error {
			err := s.ensureBootstrapped(r.Context())
			if err != nil {
				return err
			}
			return nil
		},
		Handler:       h,
		FailedHandler: handleConfigError,
	}
}

func (s *server) requireKubeconfig(h http.HandlerFunc) http.Handler {
	return &handler.PreconditionHandler{
		Precondition:  func(r *http.Request) error { return s.checkKubeconfig() },
		Handler:       h,
		FailedHandler: handleConfigError,
	}
}

func handleConfigError(w http.ResponseWriter, r *http.Request, err error) {
	if sce, ok := err.(ServerConfigError); ok {
		if sce.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
		} else {
			http.Redirect(w, r, "/support", http.StatusFound)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func defaultKubeconfigExists() bool {
	if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil {
		return true
	} else {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "could not check kubeconfig file: %v", err)
		}
		return false
	}
}

func (s *server) initPackageStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.pkgClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var pkgList v1alpha1.ClusterPackageList
				err := pkgClient.ClusterPackages().GetAll(ctx, &pkgList)
				return &pkgList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.ClusterPackages().Watch(ctx, withDecreasedTimeout(options))
			},
		},
		&v1alpha1.ClusterPackage{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.ClusterPackage); ok {
					s.broadcastRefreshTriggers(pkg)
				}
			},
			UpdateFunc: func(oldObj, newObj any) {
				if pkg, ok := newObj.(*v1alpha1.ClusterPackage); ok {
					s.broadcastRefreshTriggers(pkg)
				}
			},
			DeleteFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.ClusterPackage); ok {
					s.broadcastRefreshTriggers(pkg)
				}
			},
		},
	)
}

func (s *server) broadcastRefreshTriggers(pkg *v1alpha1.ClusterPackage) {
	s.sseHub.Broadcast <- &sse{
		event: "refresh-pkg-overview",
	}
	s.sseHub.Broadcast <- &sse{
		event: fmt.Sprintf("refresh-pkg-detail-%s", pkg.Name),
	}
}

func (s *server) initPackageInfoStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.pkgClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var packageInfoList v1alpha1.PackageInfoList
				err := pkgClient.PackageInfos().GetAll(ctx, &packageInfoList)
				return &packageInfoList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.PackageInfos().Watch(ctx, withDecreasedTimeout(options))
			},
		},
		&v1alpha1.PackageInfo{},
		0,
		cache.ResourceEventHandlerFuncs{},
	)
}

func (s *server) initPackageRepoStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.pkgClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var repositoryList v1alpha1.PackageRepositoryList
				err := pkgClient.PackageRepositories().GetAll(ctx, &repositoryList)
				return &repositoryList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.PackageRepositories().Watch(ctx, withDecreasedTimeout(options))
			},
		},
		&v1alpha1.PackageRepository{},
		0,
		cache.ResourceEventHandlerFuncs{}, // TODO we might also want to update here?
	)
}

func withDecreasedTimeout(opts metav1.ListOptions) metav1.ListOptions {
	if opts.TimeoutSeconds != nil {
		// reuse the randomized timeout and divide by 10
		t := *opts.TimeoutSeconds / 10
		opts.TimeoutSeconds = &t
	}
	return opts
}

func (s *server) isUpdateAvailable(ctx context.Context, packages ...string) bool {
	if tx, err := update.NewUpdater(ctx).Prepare(ctx, packages); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		return false
	} else {
		return !tx.IsEmpty()
	}
}

func (s *server) respondAlertAndLog(w http.ResponseWriter, err error, wrappingMsg string, alertType string) {
	if wrappingMsg != "" {
		err = fmt.Errorf("%v: %w", wrappingMsg, err)
	}
	fmt.Fprintf(os.Stderr, "%v\n", err)
	s.respondAlert(w, err.Error(), alertType)
}

func (s *server) respondAlert(w http.ResponseWriter, message string, alertType string) {
	w.Header().Add("Hx-Reselect", "div.alert") // overwrite any existing hx-select (which was a little intransparent sometimes)
	w.Header().Add("Hx-Reswap", "afterbegin")
	w.WriteHeader(http.StatusBadRequest)
	err := s.templates.alertTmpl.Execute(w, map[string]any{
		"Message":     message,
		"Dismissible": true,
		"Type":        alertType,
	})
	checkTmplError(err, "alert")
}

func isPortConflictError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if osErr, ok := opErr.Err.(*os.SyscallError); ok {
			return osErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}

func (s *server) repositoryconfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Handle GET request: Render the repository.html template
		s.handleGetRepositoryConfig(w, r)
	case http.MethodPost:
		// Handle POST request: Update the repository configuration
		// s.handlePostRepositoryConfig(w, r)
		fmt.Printf("POST Mehtod Part")
	default:
		// Method not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) handleGetRepositoryConfig(w http.ResponseWriter, r *http.Request) {
	reposName = mux.Vars(r)["repoName"]
	repos := v1alpha1.PackageRepositoryList{}
	if err := s.pkgClient.PackageRepositories().Get(r.Context(), reposName, repos); err != nil {
		s.respondAlertAndLog(w, err, "Failed to fetch repositories", "danger")
		return
	}

	tmplErr := s.templates.repositoryTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Repositories": repos.Items,
	}, nil))
	checkTmplError(tmplErr, "repository")
}

func handlePostRepositoryConfig(w http.ResponseWriter, r *http.Request){

}