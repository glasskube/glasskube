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
	"time"

	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"

	"github.com/glasskube/glasskube/internal/web/util"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/repo"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	repotypes "github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/telemetry"
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
	nonCachedClient    client.PackageV1Alpha1Client
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
	router.Use(telemetry.HttpMiddleware(telemetry.WithPathRedactor(packagesPathRedactor)))
	router.PathPrefix("/static/").Handler(fileServer)
	router.Handle("/favicon.ico", fileServer)
	router.HandleFunc("/events", s.sseHub.handler)
	router.HandleFunc("/support", s.supportPage)
	router.HandleFunc("/kubeconfig", s.kubeconfigPage)
	router.Handle("/bootstrap", s.requireKubeconfig(s.bootstrapPage))
	router.Handle("/kubeconfig/persist", s.requireKubeconfig(s.persistKubeconfig))
	// overview pages
	router.Handle("/packages", s.requireReady(s.packages))
	router.Handle("/clusterpackages", s.requireReady(s.clusterPackages))

	// detail page endpoints
	pkgBasePath := "/packages/{manifestName}"
	installedPkgBasePath := pkgBasePath + "/{namespace}/{name}"
	clpkgBasePath := "/clusterpackages/{pkgName}"
	router.Handle(pkgBasePath, s.requireReady(s.packageDetail))
	router.Handle(installedPkgBasePath, s.requireReady(s.packageDetail))
	router.Handle(clpkgBasePath, s.requireReady(s.clusterPackageDetail))
	// discussion endpoints
	router.Handle(pkgBasePath+"/discussion", s.requireReady(s.packageDiscussion))
	router.Handle(installedPkgBasePath+"/discussion", s.requireReady(s.packageDiscussion))
	router.Handle(clpkgBasePath+"/discussion", s.requireReady(s.clusterPackageDiscussion))
	router.Handle(pkgBasePath+"/discussion/badge", s.requireReady(s.discussionBadge))
	router.Handle(installedPkgBasePath+"/discussion/badge", s.requireReady(s.discussionBadge))
	router.Handle(clpkgBasePath+"/discussion/badge", s.requireReady(s.discussionBadge))
	// configuration endpoints
	router.Handle(installedPkgBasePath+"/configure", s.requireReady(s.installOrConfigurePackage))
	router.Handle(clpkgBasePath+"/configure", s.requireReady(s.installOrConfigureClusterPackage))
	router.Handle(installedPkgBasePath+"/configure/advanced", s.requireReady(s.advancedPackageConfiguration))
	router.Handle(clpkgBasePath+"/configure/advanced", s.requireReady(s.advancedClusterPackageConfiguration))
	router.Handle(pkgBasePath+"/configuration/{valueName}", s.requireReady(s.packageConfigurationInput))
	router.Handle(installedPkgBasePath+"/configuration/{valueName}", s.requireReady(s.packageConfigurationInput))
	router.Handle(clpkgBasePath+"/configuration/{valueName}", s.requireReady(s.clusterPackageConfigurationInput))
	// update endpoints
	router.Handle(installedPkgBasePath+"/update", s.requireReady(s.update))
	router.Handle(clpkgBasePath+"/update", s.requireReady(s.update))
	// open endpoints
	router.Handle(installedPkgBasePath+"/open", s.requireReady(s.open))
	router.Handle(clpkgBasePath+"/open", s.requireReady(s.open))
	// uninstall endpoints
	router.Handle(installedPkgBasePath+"/uninstall", s.requireReady(s.uninstall))
	router.Handle(clpkgBasePath+"/uninstall", s.requireReady(s.uninstall))

	// configuration datalist endpoints
	router.Handle("/datalists/{valueName}/names", s.requireReady(s.namesDatalist))
	router.Handle("/datalists/{valueName}/keys", s.requireReady(s.keysDatalist))
	// settings
	router.Handle("/settings", s.requireReady(s.settingsPage))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/clusterpackages", http.StatusFound)
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

// uninstall is an endpoint, which returns the modal html for GET requests, and performs the update for POST
func (s *server) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	if r.Method == http.MethodPost {
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
	} else {
		packageHref := ""
		updates := make([]map[string]any, 0)
		updateGetters := make([]update.PackagesGetter, 0, 1)
		if pkgName != "" {
			packageHref = "/clusterpackages/" + pkgName
			// update concerns cluster packages
			if pkgName == "-" {
				// prepare updates for all installed packages
				updateGetters = append(updateGetters, update.GetAllClusterPackages())
			} else {
				// prepare update for a specific package
				updateGetters = append(updateGetters, update.GetClusterPackageWithName(pkgName))
			}
		} else {
			// update concerns namespaced packages
			packageHref = util.GetNamespacedPkgHref(manifestName, namespace, name)
			if manifestName == "-" {
				// prepare updates for all installed namespaced packages
				updateGetters = append(updateGetters, update.GetAllPackages(""))
			} else {
				// prepare update for a specific namespaced package
				updateGetters = append(updateGetters, update.GetPackageWithName(namespace, name))
			}
		}

		updater := update.NewUpdater(ctx).WithStatusWriter(statuswriter.Stderr())
		updateTx, err := updater.Prepare(ctx, updateGetters...)
		if err != nil {
			s.respondAlertAndLog(w, err, "An error occurred preparing update of "+pkgName, "danger")
			return
		}
		utId := rand.Int()
		s.updateMutex.Lock()
		s.updateTransactions[utId] = *updateTx
		s.updateMutex.Unlock()

		for _, u := range updateTx.Items {
			if u.UpdateRequired() {
				updates = append(updates, map[string]any{
					"Package":        u.Package,
					"CurrentVersion": u.Package.GetSpec().PackageInfo.Version,
					"LatestVersion":  u.Version,
				})
			}
		}
		for _, req := range updateTx.Requirements {
			updates = append(updates, map[string]any{
				"Package":        req,
				"CurrentVersion": "-",
				"LatestVersion":  req.Version,
			})
		}

		err = s.templates.pkgUpdateModalTmpl.Execute(w, map[string]any{
			"UpdateTransactionId": utId,
			"Updates":             updates,
			"PackageHref":         packageHref,
		})
		checkTmplError(err, "pkgUpdateModalTmpl")
	}
}

// uninstall is an endpoint, which returns the modal html for GET requests, and performs the uninstallation for POST
func (s *server) uninstall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	if r.Method == http.MethodPost {
		uninstaller := uninstall.NewUninstaller(s.pkgClient).WithStatusWriter(statuswriter.Stderr())
		if pkgName != "" {
			var pkg v1alpha1.ClusterPackage
			if err := s.pkgClient.ClusterPackages().Get(ctx, pkgName, &pkg); err != nil {
				s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching %v during uninstall", pkgName), "danger")
				return
			}
			if err := uninstaller.Uninstall(ctx, &pkg); err != nil {
				s.respondAlertAndLog(w, err, "An error occurred uninstalling "+pkgName, "danger")
				return
			}
		} else {
			var pkg v1alpha1.Package
			if err := s.pkgClient.Packages(namespace).Get(ctx, name, &pkg); err != nil {
				s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching %v during uninstall", name), "danger")
				return
			}
			if err := uninstaller.Uninstall(ctx, &pkg); err != nil {
				s.respondAlertAndLog(w, err, "An error occurred uninstalling "+name, "danger")
				return
			}
		}
	} else {
		if pkgName != "" {
			var pruned []string
			var err error
			// dependency checks are only necessary for clusterpackages, as there are no dependencies on namespaced packages
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
				"PackageHref": util.GetClusterPkgHref(pkgName),
			})
			checkTmplError(err, "pkgUninstallModalTmpl")
		} else {
			err := s.templates.pkgUninstallModalTmpl.Execute(w, map[string]any{
				"Namespace":   namespace,
				"Name":        name,
				"PackageHref": util.GetNamespacedPkgHref(manifestName, namespace, name),
			})
			checkTmplError(err, "pkgUninstallModalTmpl")
		}
	}
}

func (s *server) open(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	if pkgName != "" {
		var pkg v1alpha1.ClusterPackage
		if err := s.pkgClient.ClusterPackages().Get(ctx, pkgName, &pkg); err != nil {
			s.respondAlertAndLog(w, err, "Could not get ClusterPackage", "danger")
			return
		}
		s.handleOpen(ctx, w, &pkg)
	} else {
		var pkg v1alpha1.Package
		if err := s.pkgClient.Packages(namespace).Get(ctx, name, &pkg); err != nil {
			s.respondAlertAndLog(w, err, "Could not get Package", "danger")
			return
		}
		s.handleOpen(ctx, w, &pkg)
	}
}

func (s *server) handleOpen(ctx context.Context, w http.ResponseWriter, pkg ctrlpkg.Package) {
	fwName := cache.NewObjectName(pkg.GetNamespace(), pkg.GetName()).String()
	if result, ok := s.forwarders[fwName]; ok {
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		return
	}

	result, err := open.NewOpener().Open(ctx, pkg, "", 0)
	if err != nil {
		s.respondAlertAndLog(w, err, "Could not open "+pkg.GetName(), "danger")
	} else {
		s.forwarders[fwName] = result
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *server) clusterPackages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clpkgs, listErr := list.NewLister(ctx).GetClusterPackagesWithStatus(ctx, list.ListOptions{IncludePackageInfos: true})
	if listErr != nil && len(clpkgs) == 0 {
		listErr = fmt.Errorf("could not load clusterpackages: %w", listErr)
		fmt.Fprintf(os.Stderr, "%v\n", listErr)
	}

	// Call isUpdateAvailable for each installed clusterpackage.
	// This is not the same as getting all updates in a single transaction, because some dependency
	// conflicts could be resolvable by installing individual clpkgs.
	installedClpkgs := make([]ctrlpkg.Package, 0, len(clpkgs))
	clpkgUpdateAvailable := map[string]bool{}
	for _, pkg := range clpkgs {
		if pkg.ClusterPackage != nil {
			installedClpkgs = append(installedClpkgs, pkg.ClusterPackage)
		}
		clpkgUpdateAvailable[pkg.Name] = s.isUpdateAvailableForPkg(r.Context(), pkg.ClusterPackage)
	}

	overallUpdatesAvailable := false
	if len(installedClpkgs) > 0 {
		overallUpdatesAvailable = s.isUpdateAvailable(r.Context(), installedClpkgs)
	}

	tmplErr := s.templates.clusterPkgsPageTemplate.Execute(w, s.enrichPage(r, map[string]any{
		"ClusterPackages":               clpkgs,
		"ClusterPackageUpdateAvailable": clpkgUpdateAvailable,
		"UpdatesAvailable":              overallUpdatesAvailable,
		"PackageHref":                   util.GetClusterPkgHref("-"),
	}, listErr))
	checkTmplError(tmplErr, "clusterpackages")
}

func (s *server) packages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allPkgs, listErr := list.NewLister(ctx).GetPackagesWithStatus(ctx, list.ListOptions{IncludePackageInfos: true})
	if listErr != nil {
		listErr = fmt.Errorf("could not load packages: %w", listErr)
		fmt.Fprintf(os.Stderr, "%v\n", listErr)
		// TODO check again
	}

	packageUpdateAvailable := map[string]bool{}
	var installed []*list.PackagesWithStatus
	var available []*list.PackagesWithStatus
	var installedPkgs []ctrlpkg.Package
	for _, pkgsWithStatus := range allPkgs {
		if len(pkgsWithStatus.Packages) > 0 {
			for _, pkgWithStatus := range pkgsWithStatus.Packages {
				installedPkgs = append(installedPkgs, pkgWithStatus.Package)

				// Call isUpdateAvailable for each installed package.
				// This is not the same as getting all updates in a single transaction, because some dependency
				// conflicts could be resolvable by installing individual packages.
				packageUpdateAvailable[cache.MetaObjectToName(pkgWithStatus.Package).String()] =
					s.isUpdateAvailableForPkg(ctx, pkgWithStatus.Package)
			}
			installed = append(installed, pkgsWithStatus)
		} else {
			available = append(available, pkgsWithStatus)
		}
	}

	overallUpdatesAvailable := false
	if len(installedPkgs) > 0 {
		overallUpdatesAvailable = s.isUpdateAvailable(r.Context(), installedPkgs)
	}

	tmplErr := s.templates.pkgsPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"InstalledPackages":      installed,
		"AvailablePackages":      available,
		"PackageUpdateAvailable": packageUpdateAvailable,
		"UpdatesAvailable":       overallUpdatesAvailable,
		"PackageHref":            util.GetNamespacedPkgHref("-", "-", "-"),
	}, listErr))
	checkTmplError(tmplErr, "packages")
}

// installOrConfigurePackage is like installOrConfigureClusterPackage but for packages
func (s *server) installOrConfigurePackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	requestedNamespace := r.FormValue("requestedNamespace")
	requestedName := r.FormValue("requestedName")
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	enableAutoUpdate := r.FormValue("enableAutoUpdate")

	var err error
	pkg := &v1alpha1.Package{}
	var mf *v1alpha1.PackageManifest
	if err := s.pkgClient.Packages(namespace).Get(ctx, name, pkg); err != nil && !apierrors.IsNotFound(err) {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching package details of %v", name), "danger")
		return
	} else if err != nil {
		pkg = nil
	}

	repositoryName, mf, err = s.getUsedRepoAndManifest(ctx, pkg, repositoryName, manifestName, selectedVersion)
	if err != nil {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred getting manifest and repo for %s", manifestName), "danger")
		return
	}

	if values, err := extractValues(r, mf); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred parsing the form", "danger")
		return
	} else if pkg == nil {
		pkg = client.PackageBuilder(manifestName).WithVersion(selectedVersion).
			WithVersion(selectedVersion).
			WithRepositoryName(repositoryName).
			WithAutoUpdates(strings.ToLower(enableAutoUpdate) == "on").
			WithValues(values).
			WithNamespace(requestedNamespace).
			WithName(requestedName).
			BuildPackage()
		opts := metav1.CreateOptions{}
		err := install.NewInstaller(s.pkgClient).
			WithStatusWriter(statuswriter.Stderr()).
			Install(ctx, pkg, opts)
		if err != nil {
			s.respondAlertAndLog(w, err, "An error occurred installing "+manifestName, "danger")
		} else {
			s.swappingRedirect(w, "/packages", "main", "main")
			w.WriteHeader(http.StatusAccepted)
		}
	} else {
		pkg.Spec.Values = values
		if err := s.pkgClient.Packages(pkg.GetNamespace()).Update(ctx, pkg); err != nil {
			s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred updating package %v", manifestName), "danger")
			return
		}
		if _, err := s.valueResolver.Resolve(ctx, values); err != nil {
			s.respondAlertAndLog(w, err, "Some values could not be resolved: ", "warning")
		} else {
			s.respondSuccess(w)
		}
	}
}

// installOrConfigureClusterPackage is an endpoint which takes POST requests, containing all necessary parameters to either
// install a new package if it does not exist yet, or update the configuration of an existing package.
// The name of the concerned package is given in the pkgName query parameter.
// In case the given package is not installed yet in the cluster, there must be a form parameter selectedVersion
// containing which version should be installed.
// In either case, the parameters from the form are parsed and converted into ValueConfiguration objects, which are
// being set in the packages spec.
func (s *server) installOrConfigureClusterPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	enableAutoUpdate := r.FormValue("enableAutoUpdate")
	var err error
	pkg := &v1alpha1.ClusterPackage{}
	var mf *v1alpha1.PackageManifest
	if err = s.pkgClient.ClusterPackages().Get(ctx, pkgName, pkg); err != nil && !apierrors.IsNotFound(err) {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching package details of %v", pkgName), "danger")
		return
	} else if err != nil {
		pkg = nil
	}

	repositoryName, mf, err = s.getUsedRepoAndManifest(ctx, pkg, repositoryName, pkgName, selectedVersion)
	if err != nil {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred getting manifest and repo for %s", pkgName), "danger")
		return
	}

	if values, err := extractValues(r, mf); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred parsing the form", "danger")
		return
	} else if pkg == nil {
		pkg = client.PackageBuilder(pkgName).WithVersion(selectedVersion).
			WithVersion(selectedVersion).
			WithRepositoryName(repositoryName).
			WithAutoUpdates(strings.ToLower(enableAutoUpdate) == "on").
			WithValues(values).
			BuildClusterPackage()
		opts := metav1.CreateOptions{}
		err := install.NewInstaller(s.pkgClient).
			WithStatusWriter(statuswriter.Stderr()).
			Install(ctx, pkg, opts)
		if err != nil {
			s.respondAlertAndLog(w, err, "An error occurred installing "+pkgName, "danger")
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
			s.respondSuccess(w)
		}
	}
}

func (s *server) getUsedRepoAndManifest(ctx context.Context, pkg ctrlpkg.Package, repositoryName string, manifestName string, selectedVersion string) (
	string, *v1alpha1.PackageManifest, error) {

	var mf v1alpha1.PackageManifest
	if pkg.IsNil() {
		var repoClient repoclient.RepoClient
		if len(repositoryName) == 0 {
			repos, err := s.repoClientset.Meta().GetReposForPackage(manifestName)
			if err != nil {
				return "", nil, err
			}
			switch len(repos) {
			case 0:
				return "", nil, errors.New("package not found in any repository")
			case 1:
				repositoryName = repos[0].Name
				repoClient = s.repoClientset.ForRepo(repos[0])
			default:
				return "", nil, errors.New("package found in multiple repositories")
			}
		} else {
			repoClient = s.repoClientset.ForRepoWithName(repositoryName)
		}
		if err := repoClient.FetchPackageManifest(manifestName, selectedVersion, &mf); err != nil {
			return "", nil, err
		}
	} else {
		if installedMf, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err != nil {
			return "", nil, err
		} else {
			mf = *installedMf
		}
	}
	return repositoryName, &mf, nil
}

// advancedClusterPackageConfiguration is a GET+POST endpoint which can be used for advanced package installation options,
// most notably for changing the package repository and changing to a specific (maybe even lower than installed)
// version of the package.
// It is only intended to be used for already installed clusterpackages, for new clusterpackages these options exist
// anyway and should be available for every user.
func (s *server) advancedClusterPackageConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	pkg, manifest, err := describe.DescribeInstalledClusterPackage(ctx, pkgName)
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
	s.handleAdvancedConfig(ctx, &packageDetailPageContext{
		repositoryName:  repositoryName,
		selectedVersion: selectedVersion,
		manifestName:    pkgName,
		pkg:             pkg,
		manifest:        manifest,
	}, r, w)
}

// advancedPackageConfiguration is like advancedClusterPackageConfiguration but for packages
func (s *server) advancedPackageConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	repositoryName := r.FormValue("repositoryName")
	selectedVersion := r.FormValue("selectedVersion")
	pkg, manifest, err := describe.DescribeInstalledPackage(ctx, namespace, name)
	if err != nil && !apierrors.IsNotFound(err) {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching package details of installed package %v", manifestName),
			"danger")
		return
	} else if pkg == nil {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("Package %v is not installed", manifestName),
			"danger")
		return
	} else if repositoryName == "" {
		repositoryName = pkg.Spec.PackageInfo.RepositoryName
	}
	s.handleAdvancedConfig(ctx, &packageDetailPageContext{
		repositoryName:  repositoryName,
		selectedVersion: selectedVersion,
		manifestName:    manifestName,
		pkg:             pkg,
		manifest:        manifest,
	}, r, w)
}

func (s *server) handleAdvancedConfig(ctx context.Context, d *packageDetailPageContext, r *http.Request, w http.ResponseWriter) {
	var err error
	var repos []v1alpha1.PackageRepository
	if repos, err = s.repoClientset.Meta().GetReposForPackage(d.manifestName); err != nil {
		fmt.Fprintf(os.Stderr, "error getting repos for package; %v", err)
	} else if d.repositoryName == "" {
		if len(repos) == 0 {
			s.respondAlertAndLog(w, fmt.Errorf("%v not found in any repository", d.manifestName), "", "danger")
			return
		}
		for _, r := range repos {
			d.repositoryName = r.Name
			if r.IsDefaultRepository() {
				break
			}
		}
	}

	if r.Method == http.MethodGet {
		var idx repo.PackageIndex
		if err := s.repoClientset.ForRepoWithName(d.repositoryName).FetchPackageIndex(d.manifestName, &idx); err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching package index of %v in repository %v", d.manifestName, d.repositoryName),
				"danger")
			return
		}
		latestVersion := idx.LatestVersion

		if d.selectedVersion == "" {
			d.selectedVersion = latestVersion
		} else if !slices.ContainsFunc(idx.Versions, func(item repotypes.PackageIndexItem) bool {
			return item.Version == d.selectedVersion
		}) {
			d.selectedVersion = latestVersion
		}

		res, err := s.dependencyMgr.Validate(r.Context(), d.manifest, d.selectedVersion)
		if err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred validating dependencies of %v in version %v", d.manifestName, d.selectedVersion),
				"danger")
			return
		}

		err = s.templates.pkgConfigAdvancedTmpl.Execute(w, s.enrichPage(r, map[string]any{
			"Status":           client.GetStatusOrPending(d.pkg),
			"Manifest":         d.manifest,
			"LatestVersion":    latestVersion,
			"ValidationResult": res,
			"ShowConflicts":    res.Status == dependency.ValidationResultStatusConflict,
			"SelectedVersion":  d.selectedVersion,
			"PackageIndex":     &idx,
			"Repositories":     repos,
			"RepositoryName":   d.repositoryName,
			"SelfHref":         fmt.Sprintf("%s/configure/advanced", util.GetPackageHref(d.pkg, d.manifest)),
		}, err))
		checkTmplError(err, fmt.Sprintf("advanced-config (%s)", d.manifestName))
	} else if r.Method == http.MethodPost {
		d.pkg.GetSpec().PackageInfo.Version = d.selectedVersion
		if d.repositoryName != "" {
			d.pkg.GetSpec().PackageInfo.RepositoryName = d.repositoryName
		}
		switch pkg := d.pkg.(type) {
		case *v1alpha1.ClusterPackage:
			if err := s.pkgClient.ClusterPackages().Update(ctx, pkg); err != nil {
				s.respondAlertAndLog(w, err,
					fmt.Sprintf("An error occurred updating clusterpackage %v to version %v in repo %v",
						d.manifestName, d.selectedVersion, d.repositoryName),
					"danger")
				return
			} else {
				s.respondSuccess(w)
			}
		case *v1alpha1.Package:
			if err := s.pkgClient.Packages(d.pkg.GetNamespace()).Update(ctx, pkg); err != nil {
				s.respondAlertAndLog(w, err,
					fmt.Sprintf("An error occurred updating package %v to version %v in repo %v",
						d.manifestName, d.selectedVersion, d.repositoryName),
					"danger")
				return
			} else {
				s.respondSuccess(w)
			}
		default:
			panic("unexpected package type")
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
	server.nonCachedClient = client // this should never be overridden
	server.pkgClient = client       // be aware that server.pkgClient is overridden with the cached client once bootstrap check succeeded
	return nil
}

func (server *server) initWhenBootstrapped(ctx context.Context) {
	server.k8sClient = kubernetes.NewForConfigOrDie(server.restConfig)
	server.initCachedClient(context.WithoutCancel(ctx))
	server.initClientDependentComponents()
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

func (server *server) initClientDependentComponents() {
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
}

func (server *server) initCachedClient(ctx context.Context) {
	clusterPackageStore, clusterPackageController := server.initClusterPackageStoreAndController(ctx)
	packageStore, packageController := server.initPackageStoreAndController(ctx)
	packageInfoStore, packageInfoController := server.initPackageInfoStoreAndController(ctx)
	packageRepoStore, packageRepoController := server.initPackageRepoStoreAndController(ctx)
	server.pkgClient = server.nonCachedClient.WithStores(clusterPackageStore, packageStore, packageInfoStore, packageRepoStore)

	clpkgVerifier := newVerifier(server.restConfig, clusterPackageVerifyLister)
	pkgVerifier := newVerifier(server.restConfig, packageVerifyLister)
	pkgInfoVerifier := newVerifier(server.restConfig, packageInfoVerifyLister)
	pkgRepoVerifier := newVerifier(server.restConfig, packageRepoVerifyLister)

	go clusterPackageController.Run(ctx.Done())
	go packageController.Run(ctx.Done())
	go packageInfoController.Run(ctx.Done())
	go packageRepoController.Run(ctx.Done())

	go server.broadcastUpdatesWhenInitiallySynced(clusterPackageController, packageController, packageInfoController, packageRepoController)

	go func() {
		for {
			select {
			case err := <-clpkgVerifier.errCh:
				server.handleVerificationError(err)
			case err := <-pkgVerifier.errCh:
				server.handleVerificationError(err)
			case err := <-pkgInfoVerifier.errCh:
				server.handleVerificationError(err)
			case err := <-pkgRepoVerifier.errCh:
				server.handleVerificationError(err)
			}
			cliutils.ExitWithError()
		}
	}()

	go clpkgVerifier.start(ctx, server.pkgClient, 5)
	go pkgVerifier.start(ctx, server.pkgClient, 10)
	go pkgInfoVerifier.start(ctx, server.pkgClient, 10)
	go pkgRepoVerifier.start(ctx, server.pkgClient, 30)
}

func (s *server) broadcastUpdatesWhenInitiallySynced(controllers ...cache.Controller) {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		if s.allControllersInitiallySynced(controllers...) {
			s.sseHub.Broadcast <- &sse{
				event: refreshClusterPkgOverview,
			}
			// TODO all possible refreshes
			break
		}
		<-tick.C
	}
}

func (s *server) allControllersInitiallySynced(controllers ...cache.Controller) bool {
	for _, c := range controllers {
		if !c.HasSynced() {
			return false
		}
	}
	return true
}

func (s *server) handleVerificationError(err error) {
	fmt.Fprintf(os.Stderr, "\n\n\n\nOUT OF SYNC – Local cache is probably outdated: %v\n", err)
	fmt.Fprintf(os.Stderr, "This is a known issue, see https://github.com/glasskube/glasskube/issues/838 – "+
		"As a consequence, the UI will appear stuck.\n")
	fmt.Fprintf(os.Stderr, "The server will stop now, please restart it manually and reload the UI in the browser! "+
		"(Of course we will fix this, sorry.)\n\n\n\n\n")
	telemetry.ReportCacheVerificationError(err)
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

func (s *server) initClusterPackageStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.nonCachedClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var pkgList v1alpha1.ClusterPackageList
				err := pkgClient.ClusterPackages().GetAll(ctx, &pkgList)
				return &pkgList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.ClusterPackages().Watch(ctx, options)
			},
		},
		&v1alpha1.ClusterPackage{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.ClusterPackage); ok {
					s.broadcastClusterPackageRefreshTriggers(pkg)
				}
			},
			UpdateFunc: func(oldObj, newObj any) {
				if pkg, ok := newObj.(*v1alpha1.ClusterPackage); ok {
					s.broadcastClusterPackageRefreshTriggers(pkg)
				}
			},
			DeleteFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.ClusterPackage); ok {
					s.broadcastClusterPackageRefreshTriggers(pkg)
					fwName := pkg.GetName()
					if result, ok := s.forwarders[fwName]; ok {
						result.Stop()
						delete(s.forwarders, fwName)
					}
				}
			},
		},
	)
}

func (s *server) initPackageStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.nonCachedClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var pkgList v1alpha1.PackageList
				err := pkgClient.Packages("").GetAll(ctx, &pkgList)
				return &pkgList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.Packages("").Watch(ctx, options)
			},
		},
		&v1alpha1.Package{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcastPackageRefreshTriggers(pkg)
				}
			},
			UpdateFunc: func(oldObj, newObj any) {
				if pkg, ok := newObj.(*v1alpha1.Package); ok {
					s.broadcastPackageRefreshTriggers(pkg)
				}
			},
			DeleteFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcastPackageRefreshTriggers(pkg)
					fwName := cache.ObjectName{Namespace: pkg.GetNamespace(), Name: pkg.GetName()}.String()
					if result, ok := s.forwarders[fwName]; ok {
						result.Stop()
						delete(s.forwarders, fwName)
					}
				}
			},
		},
	)
}

func (s *server) broadcastClusterPackageRefreshTriggers(pkg *v1alpha1.ClusterPackage) {
	s.sseHub.Broadcast <- &sse{
		event: refreshClusterPkgOverview,
	}
	s.sseHub.Broadcast <- &sse{
		event: fmt.Sprintf("refresh-pkg-detail-%s", pkg.Name),
	}
}

func (s *server) broadcastPackageRefreshTriggers(pkg *v1alpha1.Package) {
	s.sseHub.Broadcast <- &sse{
		event: refreshPkgOverview,
	}
	s.sseHub.Broadcast <- &sse{
		event: fmt.Sprintf("refresh-pkg-detail-%s-%s", pkg.Namespace, pkg.Name),
	}
}

func (s *server) initPackageInfoStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.nonCachedClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var packageInfoList v1alpha1.PackageInfoList
				err := pkgClient.PackageInfos().GetAll(ctx, &packageInfoList)
				return &packageInfoList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.PackageInfos().Watch(ctx, options)
			},
		},
		&v1alpha1.PackageInfo{},
		0,
		cache.ResourceEventHandlerFuncs{},
	)
}

func (s *server) initPackageRepoStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.nonCachedClient
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var repositoryList v1alpha1.PackageRepositoryList
				err := pkgClient.PackageRepositories().GetAll(ctx, &repositoryList)
				return &repositoryList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.PackageRepositories().Watch(ctx, options)
			},
		},
		&v1alpha1.PackageRepository{},
		0,
		cache.ResourceEventHandlerFuncs{}, // TODO we might also want to update here?
	)
}

func (s *server) isUpdateAvailableForPkg(ctx context.Context, pkg ctrlpkg.Package) bool {
	if pkg.IsNil() {
		return false
	}
	return s.isUpdateAvailable(ctx, []ctrlpkg.Package{pkg})
}

func (s *server) isUpdateAvailable(ctx context.Context, pkgs []ctrlpkg.Package) bool {
	if tx, err := update.NewUpdater(ctx).Prepare(ctx, update.GetExact(pkgs)); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		return false
	} else {
		return !tx.IsEmpty()
	}
}

func isPortConflictError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if osErr, ok := opErr.Err.(*os.SyscallError); ok {
			return osErr.Err == syscall.EADDRINUSE
		}
	}
	return false
}
