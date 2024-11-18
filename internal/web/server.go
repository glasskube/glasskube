package web

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/glasskube/glasskube/internal/web/handlers"
	webopen "github.com/glasskube/glasskube/internal/web/open"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/internal/web/types"

	"github.com/glasskube/glasskube/internal/web/sse"
	"github.com/glasskube/glasskube/internal/web/sse/refresh"

	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/web/middleware"
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
	coreListers             *types.CoreListers
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

func (s *server) CoreListers() *types.CoreListers {
	return s.coreListers
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

	responder.Init(webFs)
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
	router.Handle("GET /settings", s.requireReady(handlers.GetSettings))
	router.Handle("POST /settings", s.requireReady(handlers.PostSettings))
	router.Handle("GET /settings/repository/{repoName}", s.requireReady(handlers.GetRepository))
	router.Handle("POST /settings/repository/{repoName}", s.requireReady(handlers.PostRepository))

	// overview
	router.Handle("GET /clusterpackages", s.requireReady(handlers.GetClusterPackages))
	router.Handle("GET /packages", s.requireReady(handlers.GetPackages))

	// package detail
	router.Handle("GET /clusterpackages/{manifestName}", s.requireReady(handlers.GetClusterPackageDetail))
	router.Handle("GET /packages/{manifestName}", s.requireReady(handlers.GetPackageDetail))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}", s.requireReady(handlers.GetPackageDetail))

	// installation/update + configuration
	router.Handle("POST /clusterpackages/{manifestName}", s.requireReady(handlers.PostClusterPackageDetail))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}", s.requireReady(handlers.PostPackageDetail))

	// discussion
	router.Handle("POST /giscus", s.requireReady(handlers.PostGiscus))
	router.Handle("GET /clusterpackages/{manifestName}/discussion", s.requireReady(handlers.GetClusterPackageDiscussion))
	router.Handle("GET /packages/{manifestName}/discussion", s.requireReady(handlers.GetPackageDiscussion))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/discussion", s.requireReady(handlers.GetPackageDiscussion))
	router.Handle("GET /clusterpackages/{manifestName}/discussion/badge", s.requireReady(handlers.GetDiscussionBadge))
	router.Handle("GET /packages/{manifestName}/discussion/badge", s.requireReady(handlers.GetDiscussionBadge))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/discussion/badge", s.requireReady(handlers.GetDiscussionBadge))

	// configuration
	router.Handle("GET /clusterpackages/{manifestName}/configuration/{valueName}", s.requireReady(handlers.GetClusterPackageConfigurationInput))
	router.Handle("GET /packages/{manifestName}/configuration/{valueName}", s.requireReady(handlers.GetPackageConfigurationInput))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/configuration/{valueName}", s.requireReady(handlers.GetPackageConfigurationInput))

	// datalists
	router.Handle("GET /datalists/{valueName}/names", s.requireReady(handlers.GetNamesDatalist))
	router.Handle("GET /datalists/{valueName}/keys", s.requireReady(handlers.GetKeysDatalist))

	// open
	router.Handle("POST /clusterpackages/{manifestName}/open", s.requireReady(handlers.PostOpenClusterPackage))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}/open", s.requireReady(handlers.PostOpenPackage))

	// uninstall
	router.Handle("GET /clusterpackages/{manifestName}/uninstall", s.requireReady(handlers.GetUninstallClusterPackage))
	router.Handle("POST /clusterpackages/{manifestName}/uninstall", s.requireReady(handlers.PostUninstallClusterPackage))
	router.Handle("GET /packages/{manifestName}/{namespace}/{name}/uninstall", s.requireReady(handlers.GetUninstallPackage))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}/uninstall", s.requireReady(handlers.PostUninstallPackage))

	// suspend
	router.Handle("POST /clusterpackages/{manifestName}/suspend", s.requireReady(handlers.PostSuspend))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}/suspend", s.requireReady(handlers.PostSuspend))
	router.Handle("POST /clusterpackages/{manifestName}/resume", s.requireReady(handlers.PostResume))
	router.Handle("POST /packages/{manifestName}/{namespace}/{name}/resume", s.requireReady(handlers.PostResume))

	// setup
	router.HandleFunc("GET /support", s.supportPage)
	router.HandleFunc("GET /kubeconfig", s.getKubeconfigPage)
	router.HandleFunc("POST /kubeconfig", s.postKubeconfig)
	router.Handle("GET /bootstrap", s.requireKubeconfig(s.getBootstrap))
	router.Handle("POST /bootstrap", s.requireKubeconfig(s.postBootstrap))
	router.Handle("POST /kubeconfig/persist", s.requireKubeconfig(s.persistKubeconfig))

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/clusterpackages", http.StatusFound)
	})
	telemetryMiddleware := telemetry.HttpMiddleware(telemetry.WithPathRedactor(packagesPathRedactor))
	http.Handle("/", telemetryMiddleware(s.enrichContext(router)))

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

type supportPageData struct {
	types.TemplateContextHolder
	KubeconfigDefaultLocation string
	Err                       error
}

func (s *server) supportPage(w http.ResponseWriter, r *http.Request) {
	if err := s.ensureBootstrapped(r.Context()); err != nil {
		if err.BootstrapMissing() {
			http.Redirect(w, r, "/bootstrap", http.StatusFound)
			return
		}
		responder.SendPage(w, r, "pages/support", responder.ContextualizedTemplate(&supportPageData{
			KubeconfigDefaultLocation: clientcmd.RecommendedHomeFile,
			Err:                       err,
		}))
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

type bootstrapPageData struct {
	types.TemplateContextHolder
	Err error
}

func (s *server) getBootstrap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	isBootstrapped, err := bootstrap.IsBootstrapped(ctx, s.restConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nFailed to check whether Glasskube is bootstrapped: %v\n\n", err)
	} else if isBootstrapped {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	responder.SendPage(w, r, "pages/bootstrap", responder.ContextualizedTemplate(&bootstrapPageData{
		Err: err,
	}))
}

func (s *server) postBootstrap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client := bootstrap.NewBootstrapClient(s.restConfig)
	if _, err := client.Bootstrap(ctx, bootstrap.DefaultOptions()); err != nil {
		fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
		responder.SendComponent(w, r, "components/bootstrap-failure")
	} else {
		responder.SendComponent(w, r, "components/bootstrap-success")
	}
}

func (s *server) postKubeconfig(w http.ResponseWriter, r *http.Request) {
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

	s.getKubeconfigPage(w, r)
}

type kubeconfigPageData struct {
	types.TemplateContextHolder
	ConfigErr                 error
	KubeconfigDefaultLocation string
	DefaultKubeconfigExists   bool
}

func (s *server) getKubeconfigPage(w http.ResponseWriter, r *http.Request) {
	configErr := s.checkKubeconfig()
	responder.SendPage(w, r, "pages/kubeconfig", responder.ContextualizedTemplate(&kubeconfigPageData{
		ConfigErr:                 configErr,
		KubeconfigDefaultLocation: clientcmd.RecommendedHomeFile,
		DefaultKubeconfigExists:   defaultKubeconfigExists(),
	}))
}

func (s *server) persistKubeconfig(w http.ResponseWriter, r *http.Request) {
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
}

func (server *server) loadBytesConfig(data []byte) {
	server.configLoader = &bytesConfigLoader{data}
}

func (server *server) checkKubeconfig() ServerConfigError {
	if server.pkgClient == nil {
		return server.initKubeConfigAndStartListers()
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

func (server *server) initKubeConfigAndStartListers() ServerConfigError {
	restConfig, rawConfig, err := server.LoadConfig()
	if err != nil {
		return newKubeconfigErr(err)
	}
	client, err := client.New(restConfig)
	if err != nil {
		return newKubeconfigErr(err)
	}

	server.restConfig = restConfig
	server.rawConfig = rawConfig
	server.nonCachedClient = client // this should never be overridden
	server.pkgClient = client       // be aware that server.pkgClient is overridden with the cached client once bootstrap check succeeded

	server.k8sClient = kubernetes.NewForConfigOrDie(server.restConfig)
	factory := informers.NewSharedInformerFactory(server.k8sClient, 0)
	c := make(chan struct{})
	namespaceLister := factory.Core().V1().Namespaces().Lister()
	configMapLister := factory.Core().V1().ConfigMaps().Lister()
	secretLister := factory.Core().V1().Secrets().Lister()
	deploymentLister := factory.Apps().V1().Deployments().Lister()
	server.coreListers = &types.CoreListers{
		NamespaceLister:  &namespaceLister,
		ConfigMapLister:  &configMapLister,
		SecretLister:     &secretLister,
		DeploymentLister: &deploymentLister,
	}
	factory.Start(c) // TODO maybe the stop channel should be something else??
	telemetry.InitClient(restConfig, &namespaceLister)
	return nil
}

func (server *server) initWhenBootstrapped(ctx context.Context) {
	server.initCachedClient(context.WithoutCancel(ctx))
	server.initClientDependentComponents()
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
	return &middleware.ContextEnrichingHandler{Source: s, Handler: h}
}

func (s *server) requireReady(h http.HandlerFunc) http.Handler {
	return &middleware.PreconditionHandler{
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
	return &middleware.PreconditionHandler{
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
