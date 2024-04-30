package web

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/glasskube/glasskube/internal/telemetry"

	"github.com/glasskube/glasskube/internal/manifestvalues"
	"k8s.io/client-go/kubernetes"

	"github.com/Masterminds/semver/v3"

	"github.com/glasskube/glasskube/internal/web/components/pkg_config_input"

	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/dependency"

	"github.com/glasskube/glasskube/pkg/update"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
	"github.com/glasskube/glasskube/internal/web/components/pkg_update_alert"
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
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
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

const TriggerRefreshPackageDetail = "gk:refresh-package-detail"

type ServerOptions struct {
	Host       string
	Port       int32
	Kubeconfig string
}

func NewServer(options ServerOptions) *server {
	server := server{
		ServerOptions:      options,
		configLoader:       &defaultConfigLoader{options.Kubeconfig},
		forwarders:         make(map[string]*open.OpenResult),
		updateTransactions: make(map[int]update.UpdateTransaction),
	}
	return &server
}

type server struct {
	ServerOptions
	configLoader
	listener              net.Listener
	restConfig            *rest.Config
	rawConfig             *api.Config
	pkgClient             client.PackageV1Alpha1Client
	wsHub                 *WsHub
	packageStore          cache.Store
	packageController     cache.Controller
	packageInfoStore      cache.Store
	packageInfoController cache.Controller
	forwarders            map[string]*open.OpenResult
	dependencyMgr         *dependency.DependendcyManager
	updateMutex           sync.Mutex
	updateTransactions    map[int]update.UpdateTransaction
	valueResolver         *manifestvalues.Resolver
	isBootstrapped        bool
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

func (s *server) broadcastPkg(pkg *v1alpha1.Package, status *client.PackageStatus, installedManifest *v1alpha1.PackageManifest) {
	go func() {
		var bf bytes.Buffer
		err := pkg_update_alert.Render(&bf, pkgUpdateAlertTmpl, s.isUpdateAvailable(context.TODO()))
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_update_alert.TemplateId, pkg.Name))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()

	updateAvailable := false
	if installedManifest != nil {
		updateAvailable = s.isUpdateAvailable(context.TODO(), pkg.Name)
	}

	go func() {
		var bf bytes.Buffer
		err := pkg_overview_btn.Render(&bf, pkgOverviewBtnTmpl, pkg, status, installedManifest, updateAvailable)
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_overview_btn.TemplateId, pkg.Name))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()

	go func() {
		var bf bytes.Buffer
		err := pkg_detail_btns.Render(&bf, pkgDetailBtnsTmpl, pkg, status, installedManifest, updateAvailable)
		checkTmplError(err, fmt.Sprintf("%s (%s)", pkg_detail_btns.TemplateId, pkg.Name))
		if err == nil {
			s.wsHub.Broadcast <- bf.Bytes()
		}
	}()
}

func (s *server) Start(ctx context.Context) error {
	if s.listener != nil {
		return errors.New("server is already listening")
	}

	parseTemplates()
	if config.IsDevBuild() {
		if err := watchTemplates(); err != nil {
			fmt.Fprintf(os.Stderr, "templates will not be parsed after changes: %v\n", err)
		}
	}
	s.wsHub = NewHub()
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
	router.HandleFunc("/ws", s.wsHub.handler)
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
	router.Handle("/packages/{pkgName}/configure", s.requireReady(s.installOrConfigurePackage))
	router.Handle("/packages/{pkgName}/configuration/{valueName}", s.requireReady(s.packageConfigurationInput))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/packages", http.StatusFound) })
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
	_ = cliutils.OpenInBrowser("http://" + bindAddr)

	go s.wsHub.Run()
	server := &http.Server{}
	err = server.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *server) updateModal(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	pkgs := make([]string, 0, 1)
	if pkgName != "" {
		pkgs = append(pkgs, pkgName)
	}

	updates := make([]map[string]any, 0)
	updater := update.NewUpdater(s.pkgClient).WithStatusWriter(statuswriter.Stderr())
	ut, err := updater.Prepare(r.Context(), pkgs)
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

	err = pkgUpdateModalTmpl.Execute(w, map[string]any{
		"UpdateTransactionId": utId,
		"Updates":             updates,
		"PackageName":         pkgName,
	})
	checkTmplError(err, "pkgUpdateModalTmpl")
}

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	updater := update.NewUpdater(s.pkgClient).WithStatusWriter(statuswriter.Stderr())
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()
	utIdStr := r.FormValue("updateTransactionId")
	if utId, err := strconv.Atoi(utIdStr); err != nil {
		s.respondAlertAndLog(w, err, "Failed to parse updateTransactionId", "danger")
		return
	} else if ut, ok := s.updateTransactions[utId]; !ok {
		s.respondAlert(w, fmt.Sprintf("Failed to find UpdateTransaction with ID %d", utId), "danger")
		return
	} else if err = updater.Apply(ctx, &ut); err != nil {
		delete(s.updateTransactions, utId)
		s.respondAlertAndLog(w, err, "An error occurred during the update", "danger")
		return
	} else {
		delete(s.updateTransactions, utId)
		addHxTrigger(w, TriggerRefreshPackageDetail)
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
	err = pkgUninstallModalTmpl.Execute(w, map[string]any{
		"PackageName": pkgName,
		"Pruned":      pruned,
		"Err":         err,
	})
	checkTmplError(err, "pkgUninstallModalTmpl")
}

func (s *server) uninstall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgName := r.FormValue("packageName")
	var pkg v1alpha1.Package
	if err := s.pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching %v during uninstall", pkgName), "danger")
		return
	}
	if err := uninstall.NewUninstaller(s.pkgClient).
		WithStatusWriter(statuswriter.Stderr()).
		Uninstall(ctx, &pkg); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred uninstalling "+pkgName, "danger")
		return
	}
	addHxTrigger(w, TriggerRefreshPackageDetail)
}

func (s *server) open(w http.ResponseWriter, r *http.Request) {
	pkgName := r.FormValue("packageName")
	if result, ok := s.forwarders[pkgName]; ok {
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		return
	}

	result, err := open.NewOpener().Open(r.Context(), pkgName, "", 0)
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
	packages, err := list.GetPackagesWithStatus(s.pkgClient, r.Context(), list.ListOptions{IncludePackageInfos: true})
	if err != nil {
		err = fmt.Errorf("could not load packages: %w\n", err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	// Call isUpdateAvailable for each installed package.
	// This is not the same as getting all updates in a single transaction, because some dependency
	// conflicts could be resolvable by installing individual packages.
	packageUpdateAvailable := map[string]bool{}
	for _, pkg := range packages {
		packageUpdateAvailable[pkg.Name] = pkg.Package != nil && s.isUpdateAvailable(r.Context(), pkg.Name)
	}

	err = pkgsPageTmpl.Execute(w, s.enrichWithErrorAndWarnings(r.Context(), map[string]any{
		"Packages":               packages,
		"PackageUpdateAvailable": packageUpdateAvailable,
		"UpdatesAvailable":       s.isUpdateAvailable(r.Context()),
	}, err))
	checkTmplError(err, "packages")
}

func (s *server) packageDetail(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkgName"]
	selectedVersion := r.FormValue("selectedVersion")
	pkg, status, manifest, _, err := describe.DescribePackage(r.Context(), pkgName)
	autoUpdate := clientutils.AutoUpdateString(pkg, "Disabled")
	if err != nil {
		err = fmt.Errorf("An error occurred fetching package details of %v: %w\n", pkgName, err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	var idx repo.PackageIndex
	if err := repo.FetchPackageIndex("", pkgName, &idx); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred fetching versions of "+pkgName, "danger")
		return
	}
	if selectedVersion == "" {
		selectedVersion = idx.LatestVersion
	}
	if selectedVersion != idx.LatestVersion {
		var mf v1alpha1.PackageManifest
		if err := repo.FetchPackageManifest("", pkgName, selectedVersion, &mf); err != nil {
			s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching manifest of %v in version %v", pkgName, selectedVersion), "danger")
			return
		}
		manifest = &mf
	}

	res, err := s.dependencyMgr.Validate(r.Context(), manifest, selectedVersion)
	if err != nil {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred validating dependencies of %v in version %v", pkgName, selectedVersion), "danger")
		return
	}
	err = pkgPageTmpl.Execute(w, s.enrichWithErrorAndWarnings(r.Context(), map[string]any{
		"Package":           pkg,
		"Status":            status,
		"Manifest":          manifest,
		"LatestVersion":     idx.LatestVersion,
		"UpdateAvailable":   pkg != nil && s.isUpdateAvailable(r.Context(), pkgName),
		"AutoUpdate":        autoUpdate,
		"ValidationResult":  res,
		"ShowConflicts":     res.Status == dependency.ValidationResultStatusConflict,
		"SelectedVersion":   selectedVersion,
		"PackageIndex":      &idx,
		"ShowConfiguration": (pkg != nil && len(manifest.ValueDefinitions) > 0 && pkg.DeletionTimestamp.IsZero()) || pkg == nil,
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
	pkgName := mux.Vars(r)["pkgName"]
	selectedVersion := r.FormValue("selectedVersion")
	enableAutoUpdate := r.FormValue("enableAutoUpdate")
	pkg, _, manifest, _, err := describe.DescribePackage(r.Context(), pkgName)
	if err != nil {
		s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching package details of %v", pkgName), "danger")
		return
	} else if pkg == nil {
		var mf v1alpha1.PackageManifest
		if err := repo.FetchPackageManifest("", pkgName, selectedVersion, &mf); err != nil {
			s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred fetching manifest of %v in version %v", pkgName, selectedVersion), "danger")
			return
		}
		manifest = &mf
	}

	if values, err := extractValues(r, manifest); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred parsing the form", "danger")
		return
	} else if pkg == nil {
		pkg = client.PackageBuilder(pkgName).
			WithVersion(selectedVersion).
			WithAutoUpdates(strings.ToLower(enableAutoUpdate) == "on").
			WithValues(values).
			Build()
		err := install.NewInstaller(s.pkgClient).
			WithStatusWriter(statuswriter.Stderr()).
			Install(r.Context(), pkg)
		if err != nil {
			s.respondAlertAndLog(w, err, "An error occurred installing "+pkgName, "danger")
			return
		}
		addHxTrigger(w, TriggerRefreshPackageDetail)
	} else {
		pkg.Spec.Values = values
		if err := s.pkgClient.Packages().Update(r.Context(), pkg); err != nil {
			s.respondAlertAndLog(w, err, fmt.Sprintf("An error occurred updating package %v", pkgName), "danger")
			return
		}
		if _, err := s.valueResolver.Resolve(r.Context(), values); err != nil {
			s.respondAlertAndLog(w, err, "Some values could not be resolved: ", "warning")
		} else {
			err := alertTmpl.Execute(w, map[string]any{
				"Message":     "Configuration updated successfully",
				"Dismissible": true,
				"Type":        "success",
			})
			checkTmplError(err, "success")
		}
	}
}

// packageConfigurationInput is a GET endpoint, which returns an html snippet containing an input container.
// The endpoint requires the pkgName query parameter to be set, as well as the valueName query parameter (which holds
// the name of the desired value according to the package value definitions).
// An optional query parameter refKind can be passed to request the snippet in a certain variant, where the accepted
// refKind values are: ConfigMap, Secret, Package. If no refKind is given, the "regular" input is returned.
// In any case, the input container consists of a button where the user can change the type of reference or remove the
// reference, and the actual input field(s).
func (s *server) packageConfigurationInput(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkgName"]
	pkg, _, manifest, _, err := describe.DescribePackage(r.Context(), pkgName)
	if err != nil {
		err = fmt.Errorf("An error occurred fetching package details of %v: %w\n", pkgName, err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	valueName := mux.Vars(r)["valueName"]
	refKind := r.URL.Query().Get("refKind")
	if valueDefinition, ok := manifest.ValueDefinitions[valueName]; ok {
		input := pkg_config_input.ForPkgConfigInput(pkg, pkgName, valueName, valueDefinition, &pkg_config_input.PkgConfigInputRenderOptions{
			Autofocus:      true,
			DesiredRefKind: &refKind,
		})
		err = pkgConfigInput.Execute(w, input)
		checkTmplError(err, fmt.Sprintf("package config input (%s, %s)", pkgName, valueName))
	}
}

func (s *server) supportPage(w http.ResponseWriter, r *http.Request) {
	if err := s.ensureBootstrapped(r.Context()); err != nil {
		if err.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
			return
		}
		err := supportPageTmpl.Execute(w, &map[string]any{
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
		if err := client.Bootstrap(ctx, bootstrap.DefaultOptions()); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			err := bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-failure", nil)
			checkTmplError(err, "bootstrap-failure")
		} else {
			err := bootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-success", nil)
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
		tplErr := bootstrapPageTmpl.Execute(w, &map[string]any{
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
	tplErr := kubeconfigPageTmpl.Execute(w, map[string]any{
		"CurrentContext":            currentContext,
		"ConfigErr":                 configErr,
		"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
		"DefaultKubeconfigExists":   defaultKubeconfigExists(),
	})
	checkTmplError(tplErr, "kubeconfig")
}

func (s *server) enrichWithErrorAndWarnings(ctx context.Context, data map[string]any, err error) map[string]any {
	cacheBustStr := fmt.Sprintf("%d", rand.Int())
	data["Error"] = err
	data["CurrentContext"] = s.rawConfig.CurrentContext
	operatorVersion, clientVersion, err := s.getGlasskubeVersions(ctx)
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
		cacheBustStr = clientVersion.String()
	}
	if config.IsDevBuild() {
		data["VersionDetails"] = map[string]any{
			"OperatorVersion": config.Version,
			"ClientVersion":   config.Version,
		}
	}
	data["CacheBustingString"] = cacheBustStr
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
	k8sclient := kubernetes.NewForConfigOrDie(server.restConfig)
	server.initCachedClient(context.WithoutCancel(ctx))
	server.dependencyMgr = dependency.NewDependencyManager(clientadapter.NewPackageClientAdapter(server.pkgClient))
	server.valueResolver = manifestvalues.NewResolver(
		clientadapter.NewPackageClientAdapter(server.pkgClient),
		clientadapter.NewKubernetesClientAdapter(*k8sclient),
	)
}

func (server *server) initCachedClient(ctx context.Context) {
	server.packageStore, server.packageController = server.initPackageStoreAndController(ctx)
	server.packageInfoStore, server.packageInfoController = server.initPackageInfoStoreAndController(ctx)
	go server.packageController.Run(ctx.Done())
	go server.packageInfoController.Run(ctx.Done())
	server.pkgClient = server.pkgClient.WithStores(server.packageStore, server.packageInfoStore)
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
				var pkgList v1alpha1.PackageList
				err := pkgClient.Packages().GetAll(ctx, &pkgList)
				return &pkgList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.Packages().Watch(ctx)
			},
		},
		&v1alpha1.Package{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcastPkg(pkg, client.GetStatusOrPending(&pkg.Status), nil)
				}
			},
			UpdateFunc: func(oldObj, newObj any) {
				if pkg, ok := newObj.(*v1alpha1.Package); ok {
					ctx := client.SetupContextWithClient(ctx, s.restConfig, s.rawConfig, s.pkgClient)
					mf, err := manifest.GetInstalledManifestForPackage(ctx, *pkg)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error fetching manifest for package %v: %v\n", pkg.Name, err)
						mf = nil
					}
					s.broadcastPkg(pkg, client.GetStatusOrPending(&pkg.Status), mf)
				}
			},
			DeleteFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcastPkg(pkg, client.GetStatus(&pkg.Status), nil)
				}
			},
		},
	)
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
				return pkgClient.PackageInfos().Watch(ctx)
			},
		},
		&v1alpha1.PackageInfo{},
		0,
		cache.ResourceEventHandlerFuncs{},
	)
}

func (s *server) isUpdateAvailable(ctx context.Context, packages ...string) bool {
	if tx, err := update.NewUpdater(s.pkgClient).Prepare(ctx, packages); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		return false
	} else {
		return !tx.IsEmpty()
	}
}

func addHxTrigger(w http.ResponseWriter, trigger string) {
	w.Header().Add("HX-Trigger", trigger)
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
	err := alertTmpl.Execute(w, map[string]any{
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
