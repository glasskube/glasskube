package web

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/telemetry/annotations"
	"github.com/glasskube/glasskube/internal/web/controllers"
	webopen "github.com/glasskube/glasskube/internal/web/open"
	"github.com/glasskube/glasskube/internal/web/responder"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/glasskube/glasskube/internal/web/sse"
	"github.com/glasskube/glasskube/internal/web/sse/refresh"

	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/web/handler"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/glasskube/glasskube/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
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
		if _, err := os.Lstat(responder.BaseDir); err == nil {
			webFs = os.DirFS(responder.BaseDir)
		}
	}
}

type ServerOptions struct {
	Host               string
	Port               string
	Kubeconfig         string
	LogLevel           int
	SkipOpeningBrowser bool
}

func NewServer(options ServerOptions) *server {
	server := server{
		ServerOptions:           options,
		configLoader:            &defaultConfigLoader{options.Kubeconfig},
		stopCh:                  make(chan struct{}, 1),
		httpServerHasShutdownCh: make(chan struct{}, 1),
	}
	return &server
}

type server struct {
	ServerOptions
	configLoader
	listener                net.Listener
	restConfig              *rest.Config
	rawConfig               *api.Config
	pkgClient               client.PackageV1Alpha1Client
	nonCachedClient         client.PackageV1Alpha1Client
	repoClientset           repoclient.RepoClientset
	k8sClient               *kubernetes.Clientset
	coreListers             *clicontext.CoreListers
	broadcaster             *sse.Broadcaster
	isBootstrapped          bool
	httpServer              *http.Server
	httpServerHasShutdownCh chan struct{}
	stopCh                  chan struct{}
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

func (s *server) CoreListers() *clicontext.CoreListers {
	return s.coreListers
}

func (s *server) RepoClient() repoclient.RepoClientset {
	return s.repoClientset
}

func (s *server) GetCurrentContext() string {
	if s.rawConfig != nil {
		return s.rawConfig.CurrentContext
	}
	return ""
}

func (s *server) IsGitopsModeEnabled() bool {
	return s.isGitopsModeEnabled()
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

	responder.Init(s, webFs)
	webopen.Init(s.Host, s.stopCh)

	s.broadcaster = sse.NewBroadcaster()
	_ = s.ensureBootstrapped(ctx)

	root, err := fs.Sub(webFs, "root")
	if err != nil {
		return err
	}

	fileServer := http.FileServer(http.FS(root))

	router := http.NewServeMux()
	router.Handle("GET /static/", fileServer)
	router.Handle("GET /favicon.ico", fileServer)

	router.HandleFunc("GET /events", s.broadcaster.Handler) // TODO ??

	// settings
	router.Handle("GET /settings", s.requireReady(controllers.GetSettings))
	router.Handle("POST /settings", s.requireReady(controllers.PostSettings))
	router.Handle("GET /settings/repository/{repoName}", s.requireReady(controllers.GetRepository))
	router.Handle("POST /settings/repository/{repoName}", s.requireReady(controllers.PostRepository))

	// overview
	router.Handle("GET /clusterpackages", s.requireReady(controllers.GetClusterPackages))
	router.Handle("GET /packages", s.requireReady(controllers.GetPackages))

	// package detail
	router.Handle("GET /clusterpackages/{manifestName}", s.requireReady(controllers.GetClusterPackageDetail))
	// TODO maybe there is a more fancy way for these "duplicated" routes (also for configuration etc subpaths):
	router.Handle("GET /packages/{manifestName}", s.requireReady(controllers.GetPackageDetail))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}", s.requireReady(controllers.GetPackageDetail))

	// installation/update + configuration
	router.Handle("POST /clusterpackages/{manifestName}", s.requireReady(controllers.PostClusterPackageDetail))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}", s.requireReady(controllers.PostPackageDetail))

	// discussion
	router.Handle("POST /giscus", s.requireReady(controllers.PostGiscus))
	router.Handle("GET /clusterpackages/{manifestName}/discussion", s.requireReady(controllers.GetClusterPackageDiscussion))
	router.Handle("GET /packages/{manifestName}/discussion", s.requireReady(controllers.GetPackageDiscussion))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/discussion", s.requireReady(controllers.GetPackageDiscussion))
	router.Handle("GET /clusterpackages/{manifestName}/discussion/badge", s.requireReady(controllers.GetDiscussionBadge))
	router.Handle("GET /packages/{manifestName}/discussion/badge", s.requireReady(controllers.GetDiscussionBadge))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/discussion/badge", s.requireReady(controllers.GetDiscussionBadge))

	// configuration
	router.Handle("GET /clusterpackages/{manifestName}/configuration/{valueName}", s.requireReady(controllers.GetClusterPackageConfigurationInput))
	router.Handle("GET /packages/{manifestName}/configuration/{valueName}", s.requireReady(controllers.GetPackageConfigurationInput))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/configuration/{valueName}", s.requireReady(controllers.GetPackageConfigurationInput))

	// open
	router.Handle("POST /clusterpackages/{manifestName}/open", s.requireReady(controllers.PostOpenClusterPackage))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}/open", s.requireReady(controllers.PostOpenPackage))

	// uninstall
	router.Handle("GET /clusterpackages/{manifestName}/uninstall", s.requireReady(controllers.GetUninstallClusterPackage))
	router.Handle("POST /clusterpackages/{manifestName}/uninstall", s.requireReady(controllers.PostUninstallClusterPackage))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/uninstall", s.requireReady(controllers.GetUninstallPackage))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}/uninstall", s.requireReady(controllers.PostUninstallPackage))

	// suspend

	// datalists
	router.Handle("GET /datalists/{valueName}/names", s.requireReady(controllers.GetNamesDatalist))
	router.Handle("GET /datalists/{valueName}/keys", s.requireReady(controllers.GetKeysDatalist))

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/clusterpackages", http.StatusFound)
	})
	http.Handle("/", s.enrichContext(router))

	/*
		TODO
		router.Use(telemetry.HttpMiddleware(telemetry.WithPathRedactor(packagesPathRedactor)))

		router.HandleFunc("/support", s.supportPage)
		router.HandleFunc("/kubeconfig", s.kubeconfigPage)
		router.Handle("/bootstrap", s.requireKubeconfig(s.bootstrapPage))
		router.Handle("/kubeconfig/persist", s.requireKubeconfig(s.persistKubeconfig))

		// detail page endpoints
		pkgBasePath := "/packages/{manifestName}"
		installedPkgBasePath := pkgBasePath + "/{namespace}/{name}"
		clpkgBasePath := "/clusterpackages/{pkgName}"

		// suspend endpoints
		router.Handle(clpkgBasePath+"/suspend", s.requireReady(s.handleSuspend))
		router.Handle(clpkgBasePath+"/resume", s.requireReady(s.handleResume))
		router.Handle(installedPkgBasePath+"/suspend", s.requireReady(s.handleSuspend))
		router.Handle(installedPkgBasePath+"/resume", s.requireReady(s.handleResume))

	*/

	s.listener, err = net.Listen("tcp", net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		// if the error is "address already in use", try to get the OS to assign a random free port
		if errors.Is(err, syscall.EADDRINUSE) {
			fmt.Fprintf(os.Stderr, "could not start server: %v\n", err)
			if cliutils.YesNoPrompt("Should glasskube try to use a different (random) port?", true) {
				s.listener, err = net.Listen("tcp", net.JoinHostPort(s.Host, "0"))
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}

	browseUrl := fmt.Sprintf("http://%s", s.listener.Addr())
	fmt.Fprintln(os.Stderr, "glasskube UI is available at", browseUrl)
	if !s.SkipOpeningBrowser {
		_ = cliutils.OpenInBrowser(browseUrl)
	}

	go s.broadcaster.Run(s.stopCh)
	s.httpServer = &http.Server{}

	var receivedSig *os.Signal
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT)
		sig := <-sigint
		receivedSig = &sig
		s.shutdown()
	}()

	err = s.httpServer.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-s.httpServerHasShutdownCh
	cliutils.ExitFromSignal(receivedSig)

	return nil
}

func (s *server) shutdown() {
	close(s.stopCh)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to shutdown server: %v\n", err)
	}
	close(s.httpServerHasShutdownCh)
}

func (s *server) isGitopsModeEnabled() bool {
	if ns, err := (*s.coreListers.NamespaceLister).Get("glasskube-system"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch glasskube-system namespace: %v\n", err)
		return true
	} else {
		return annotations.IsGitopsModeEnabled(ns.Annotations)
	}
}

/*
func (s *server) supportPage(w http.ResponseWriter, r *http.Request) {
	if err := s.ensureBootstrapped(r.Context()); err != nil {
		if err.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
			return
		}
		err := templates.Templates.SupportPageTmpl.Execute(w, &map[string]any{
			"CurrentContext":            "",
			"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
			"Err":                       err,
		})
		util.CheckTmplError(err, "support")
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
			err := templates.Templates.BootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-failure", nil)
			util.CheckTmplError(err, "bootstrap-failure")
		} else {
			err := templates.Templates.BootstrapPageTmpl.ExecuteTemplate(w, "bootstrap-success", nil)
			util.CheckTmplError(err, "bootstrap-success")
		}
	} else {
		isBootstrapped, err := bootstrap.IsBootstrapped(ctx, s.restConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
		} else if isBootstrapped {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		tplErr := templates.Templates.BootstrapPageTmpl.Execute(w, &map[string]any{
			"CloudId":        telemetry.GetMachineId(),
			"CurrentContext": s.rawConfig.CurrentContext,
			"Err":            err,
		})
		util.CheckTmplError(tplErr, "bootstrap")
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
	tplErr := templates.Templates.KubeconfigPageTmpl.Execute(w, map[string]any{
		"CloudId":                   telemetry.GetMachineId(),
		"CurrentContext":            currentContext,
		"ConfigErr":                 configErr,
		"KubeconfigDefaultLocation": clientcmd.RecommendedHomeFile,
		"DefaultKubeconfigExists":   defaultKubeconfigExists(),
	})
	util.CheckTmplError(tplErr, "kubeconfig")
}

func (s *server) enrichPage(r *http.Request, data map[string]any, err error) map[string]any {
	data["CloudId"] = telemetry.GetMachineId()
	if pathParts := strings.Split(r.URL.Path, "/"); len(pathParts) >= 2 {
		data["NavbarActiveItem"] = pathParts[1]
	}
	data["Error"] = err
	data["CurrentContext"] = s.rawConfig.CurrentContext
	data["GitopsMode"] = s.isGitopsModeEnabled()
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
			"GitopsMode":          s.isGitopsModeEnabled(),
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

*/

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
	configMapLister := factory.Core().V1().ConfigMaps().Lister()
	secretLister := factory.Core().V1().Secrets().Lister()
	server.coreListers = &clicontext.CoreListers{
		NamespaceLister: &namespaceLister,
		ConfigMapLister: &configMapLister,
		SecretLister:    &secretLister,
	}
	factory.Start(c) // TODO maybe the stop channel should be something else??
}

func (server *server) initClientDependentComponents() {
	server.repoClientset = repoclient.NewClientset(
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
			server.shutdown()
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
			var allPkgs []ctrlpkg.Package

			var clpkgs v1alpha1.ClusterPackageList
			if err := s.pkgClient.ClusterPackages().GetAll(context.TODO(), &clpkgs); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get all clusterpackages to broadcast all updates: %v\n", err)
			} else {
				for _, clpkg := range clpkgs.Items {
					p := &clpkg
					allPkgs = append(allPkgs, p)
				}
			}

			var pkgs v1alpha1.PackageList
			if err := s.pkgClient.Packages("").GetAll(context.TODO(), &pkgs); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get all packages to broadcast all updates: %v\n", err)
			} else {
				for _, pkg := range pkgs.Items {
					p := &pkg
					allPkgs = append(allPkgs, p)
				}
			}

			s.broadcaster.UpdatesAvailable(refresh.RefreshTriggerAll, allPkgs...)
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

func (s *server) requireKubeconfig(h http.Handler) http.Handler {
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
	return cache.NewInformerWithOptions(cache.InformerOptions{
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var pkgList v1alpha1.ClusterPackageList
				err := pkgClient.ClusterPackages().GetAll(ctx, &pkgList)
				return &pkgList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.ClusterPackages().Watch(ctx, options)
			},
		},
		ObjectType: &v1alpha1.ClusterPackage{},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.ClusterPackage); ok {
					s.broadcaster.UpdatesAvailableForPackage(nil, pkg)
				}
			},
			UpdateFunc: func(oldObj, newObj any) {
				if oldPkg, ok := oldObj.(*v1alpha1.ClusterPackage); ok {
					if newPkg, ok := newObj.(*v1alpha1.ClusterPackage); ok {
						s.broadcaster.UpdatesAvailableForPackage(oldPkg, newPkg)
					}
				}
			},
			DeleteFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.ClusterPackage); ok {
					s.broadcaster.UpdatesAvailableForPackage(pkg, nil)
					webopen.CloseForwarders(pkg)
				}
			},
		},
	})
}

func (s *server) initPackageStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.nonCachedClient
	return cache.NewInformerWithOptions(cache.InformerOptions{
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var pkgList v1alpha1.PackageList
				err := pkgClient.Packages("").GetAll(ctx, &pkgList)
				return &pkgList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.Packages("").Watch(ctx, options)
			},
		},
		ObjectType: &v1alpha1.Package{},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcaster.UpdatesAvailableForPackage(nil, pkg)
				}
			},
			UpdateFunc: func(oldObj, newObj any) {
				if oldPkg, ok := oldObj.(*v1alpha1.Package); ok {
					if newPkg, ok := newObj.(*v1alpha1.Package); ok {
						s.broadcaster.UpdatesAvailableForPackage(oldPkg, newPkg)
					}
				}
			},
			DeleteFunc: func(obj any) {
				if pkg, ok := obj.(*v1alpha1.Package); ok {
					s.broadcaster.UpdatesAvailableForPackage(pkg, nil)
					webopen.CloseForwarders(pkg)
				}
			},
		},
	})
}

func (s *server) initPackageInfoStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.nonCachedClient
	return cache.NewInformerWithOptions(cache.InformerOptions{
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var packageInfoList v1alpha1.PackageInfoList
				err := pkgClient.PackageInfos().GetAll(ctx, &packageInfoList)
				return &packageInfoList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.PackageInfos().Watch(ctx, options)
			},
		},
		ObjectType: &v1alpha1.PackageInfo{},
		Handler:    cache.ResourceEventHandlerFuncs{},
	})
}

func (s *server) initPackageRepoStoreAndController(ctx context.Context) (cache.Store, cache.Controller) {
	pkgClient := s.nonCachedClient
	return cache.NewInformerWithOptions(cache.InformerOptions{
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				var repositoryList v1alpha1.PackageRepositoryList
				err := pkgClient.PackageRepositories().GetAll(ctx, &repositoryList)
				return &repositoryList, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return pkgClient.PackageRepositories().Watch(ctx, options)
			},
		},
		ObjectType: &v1alpha1.PackageRepository{},
		Handler:    cache.ResourceEventHandlerFuncs{}, // TODO we might also want to update here?
	})
}
