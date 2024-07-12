package telemetry

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/glasskube/glasskube/internal/config"

	"github.com/denisbrodbeck/machineid"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/telemetry/properties"
	"github.com/posthog/posthog-go"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type clientsetNamespaceGetter struct {
	client *kubernetes.Clientset
}

func (g clientsetNamespaceGetter) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	return g.client.CoreV1().Namespaces().Get(ctx, name, v1.GetOptions{})
}

type clientsetNodeLister struct {
	client *kubernetes.Clientset
}

func (l clientsetNodeLister) ListNodes(ctx context.Context) (*corev1.NodeList, error) {
	return l.client.CoreV1().Nodes().List(ctx, v1.ListOptions{})
}

type ClientTelemetry struct {
	properties properties.PropertyGetter
	machineId  string
	posthog    posthog.Client
	start      time.Time
	restConfig *rest.Config
}

var instance *ClientTelemetry

func Init() {
	instance = ForClient()
}

func InitClient(config *rest.Config) {
	instance.InitClient(config)
}

func BootstrapAttempt() {
	if instance != nil {
		command, arguments, flags := getCommandAndArgs()
		event := instance.getBaseEvent("bootstrap_attempt", true)
		event.Properties["original_command"] = command
		event.Properties["arguments"] = arguments
		event.Properties["flags"] = flags
		_ = instance.posthog.Enqueue(event)
	}
}

func BootstrapFailure(elapsed time.Duration) {
	if instance != nil {
		command, arguments, flags := getCommandAndArgs()
		event := instance.getBaseEvent("bootstrap_failure", true)
		event.Properties["original_command"] = command
		event.Properties["arguments"] = arguments
		event.Properties["flags"] = flags
		event.Properties["execution_time"] = elapsed.Milliseconds()
		_ = instance.posthog.Enqueue(event)
	}
}

func BootstrapSuccess(elapsed time.Duration) {
	if instance != nil {
		command, arguments, flags := getCommandAndArgs()
		event := instance.getBaseEvent("bootstrap_success", true)
		event.Properties["original_command"] = command
		event.Properties["arguments"] = arguments
		event.Properties["flags"] = flags
		event.Properties["execution_time"] = elapsed.Milliseconds()
		_ = instance.posthog.Enqueue(event)
	}
}

func Exit() {
	if instance != nil {
		instance.report(0, "")
		instance.close()
	}
}

func ExitFromSignal(sig os.Signal) {
	if instance != nil {
		instance.report(1, sig.String())
		instance.close()
	}
}

func ExitWithError() {
	if instance != nil {
		instance.report(1, "")
		instance.close()
	}
}

func ForClient() *ClientTelemetry {
	ct := ClientTelemetry{
		machineId: GetMachineId(),
		start:     time.Now(),
	}
	if ph, err := posthog.NewWithConfig(apiKey, posthog.Config{
		Endpoint: endpoint,
		Logger:   posthog.StdLogger(log.New(io.Discard, "", 0)),
	}); err == nil {
		ct.posthog = ph
	}
	return &ct
}

func (t *ClientTelemetry) InitClient(config *rest.Config) {
	if config != nil {
		t.restConfig = config
		if client, err := kubernetes.NewForConfig(config); err == nil {
			t.properties.NamespaceGetter = clientsetNamespaceGetter{client}
			t.properties.NodeLister = clientsetNodeLister{client}
			t.properties.DiscoveryClient = client
		}
	}
}

func GetMachineId() string {
	id, err := machineid.ProtectedID("glasskube")
	if err != nil {
		return fmt.Sprintf("fallback-id-%v-%v", rand.Int(), err.Error())
	}
	return id
}

func getCommandAndArgs() (string, []string, []string) {
	command := "-" // not allowed to be empty by posthog
	var arguments []string
	var flags []string
	if len(os.Args) >= 2 {
		command = os.Args[1]
		arguments, flags = extractFromArgs(os.Args[2:])
	}
	return command, arguments, flags
}

func extractFromArgs(args []string) ([]string, []string) {
	var arguments = make([]string, 0)
	flags := make([]string, 0)
	if len(args) > 0 {
		noMoreArguments := false
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				noMoreArguments = true
				flags = append(flags, arg)
			} else if !noMoreArguments {
				arguments = append(arguments, arg)
			}
		}
	}
	return arguments, flags
}

func (t *ClientTelemetry) getBaseEvent(event string, includeCluster bool) *posthog.Capture {
	return &posthog.Capture{
		DistinctId: t.machineId,
		Event:      event,
		Properties: properties.BuildProperties(
			properties.ForClientUser(instance.properties, includeCluster),
		),
	}
}

func (t *ClientTelemetry) report(exitCode int, reason string) {
	if !t.properties.Enabled() {
		return
	}

	operatorVersion, _ := clientutils.GetPackageOperatorVersionForConfig(t.restConfig, context.TODO())
	command, arguments, flags := getCommandAndArgs()
	duration := time.Since(t.start).Milliseconds()

	event := t.getBaseEvent(command, true)
	event.Properties["arguments"] = arguments
	event.Properties["flags"] = flags
	event.Properties["execution_time"] = duration
	event.Properties["exit_code"] = exitCode
	event.Properties["exit_reason"] = reason
	event.Properties["operator_version"] = operatorVersion

	err := t.posthog.Enqueue(event)
	if err != nil && config.IsDevBuild() {
		fmt.Fprintf(os.Stderr, "Telemetry error: %v\n", err)
	}
}

func (t *ClientTelemetry) close() {
	err := t.posthog.Close()
	if err != nil && config.IsDevBuild() {
		fmt.Fprintf(os.Stderr, "Telemetry error: %v\n", err)
	}
}

type PathRedactor = func(url string) string

type httpMiddlewareOptions struct {
	PathRedactors []PathRedactor
}

type HttpMiddlewareConfigurator = func(opt *httpMiddlewareOptions)

func WithPathRedactor(pr PathRedactor) HttpMiddlewareConfigurator {
	return func(opt *httpMiddlewareOptions) {
		opt.PathRedactors = append(opt.PathRedactors, pr)
	}
}

func HttpMiddleware(conf ...HttpMiddlewareConfigurator) func(http.Handler) http.Handler {
	var options httpMiddlewareOptions
	for _, c := range conf {
		c(&options)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			defer func() {
				if instance != nil {
					go func() {
						redactedUrl := redactedPath(r, options.PathRedactors)
						ev := instance.getBaseEvent("ui_endpoint", false)
						ev.Properties["$current_url"] = redactedUrl.String()
						ev.Properties["method"] = r.Method
						ev.Properties["path"] = redactedUrl
						ev.Properties["execution_time"] = time.Since(start).Milliseconds()
						ev.Properties["user_agent"] = r.UserAgent()
						_ = instance.posthog.Enqueue(ev)
					}()
				}
			}()
		})
	}
}

func redactedPath(r *http.Request, redactors []PathRedactor) url.URL {
	result := *r.URL
	for _, redactor := range redactors {
		result.Path = redactor(result.Path)
	}
	return result
}

func SetUserProperty(key string, value string) {
	_ = instance.posthog.Enqueue(posthog.Capture{
		DistinctId: instance.machineId,
		Event:      "$set",
		Properties: map[string]any{
			"$set": map[string]any{
				key: value,
			},
		},
	})
}

func ReportCacheVerificationError(err error) {
	ev := instance.getBaseEvent("error_cache_verification", false)
	ev.Properties["err"] = err.Error()
	_ = instance.posthog.Enqueue(ev)
}
