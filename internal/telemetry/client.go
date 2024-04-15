package telemetry

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/telemetry/properties"
	client2 "github.com/glasskube/glasskube/pkg/client"
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
		event := instance.getBaseEvent("bootstrap_attempt")
		event.Properties["original_command"] = command
		event.Properties["arguments"] = arguments
		event.Properties["flags"] = flags
		_ = instance.posthog.Enqueue(event)
	}
}

func BootstrapFailure(elapsed time.Duration) {
	if instance != nil {
		command, arguments, flags := getCommandAndArgs()
		event := instance.getBaseEvent("bootstrap_failure")
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
		event := instance.getBaseEvent("bootstrap_success")
		event.Properties["original_command"] = command
		event.Properties["arguments"] = arguments
		event.Properties["flags"] = flags
		event.Properties["execution_time"] = elapsed.Milliseconds()
		_ = instance.posthog.Enqueue(event)
	}
}

func SetupFailed() {
	if instance != nil {
		command, arguments, flags := getCommandAndArgs()
		event := instance.getBaseEvent("telemetry_setup_failed")
		event.Properties["original_command"] = command
		event.Properties["arguments"] = arguments
		event.Properties["flags"] = flags
		_ = instance.posthog.Enqueue(event)
	}
}

func Exit(ctx context.Context) {
	if instance != nil {
		instance.report(ctx, 0, "")
		instance.close()
	}
}

func ExitFromSignal(ctx context.Context, sig os.Signal) {
	if instance != nil {
		instance.report(ctx, 1, sig.String())
		instance.close()
	}
}

func ExitWithError(ctx context.Context) {
	if instance != nil {
		instance.report(ctx, 1, "")
		instance.close()
	}
}

func ForClient() *ClientTelemetry {
	ct := ClientTelemetry{
		machineId: getMachineId(),
		start:     time.Now(),
	}
	if ph, err := posthog.NewWithConfig(apiKey, posthog.Config{Endpoint: endpoint}); err == nil {
		ct.posthog = ph
	}
	return &ct
}

func (t *ClientTelemetry) InitClient(config *rest.Config) {
	if config != nil {
		if client, err := kubernetes.NewForConfig(config); err == nil {
			t.properties.NamespaceGetter = clientsetNamespaceGetter{client}
			t.properties.NodeLister = clientsetNodeLister{client}
			t.properties.DiscoveryClient = client
		}
	}
}

func getMachineId() string {
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

func (t *ClientTelemetry) getBaseEvent(event string) *posthog.Capture {
	return &posthog.Capture{
		DistinctId: t.machineId,
		Event:      event,
		Properties: properties.BuildProperties(
			properties.ForClientUser(instance.properties),
		),
	}
}

func (t *ClientTelemetry) report(ctx context.Context, exitCode int, reason string) {
	t.InitClient(client2.ConfigFromContext(ctx))

	if !t.properties.Enabled() {
		return
	}

	operatorVersion, _ := clientutils.GetPackageOperatorVersion(ctx)
	command, arguments, flags := getCommandAndArgs()
	duration := time.Since(t.start).Milliseconds()

	event := t.getBaseEvent(command)
	event.Properties["arguments"] = arguments
	event.Properties["flags"] = flags
	event.Properties["execution_time"] = duration
	event.Properties["exit_code"] = exitCode
	event.Properties["exit_reason"] = reason
	event.Properties["operator_version"] = operatorVersion

	err := t.posthog.Enqueue(event)
	if err != nil {
		// TODO only in dev?
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func (t *ClientTelemetry) close() {
	err := t.posthog.Close()
	if err != nil {
		// TODO only in dev?
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
