package telemetry

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/config"
	client2 "github.com/glasskube/glasskube/pkg/client"
	"github.com/posthog/posthog-go"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ClientTelemetry struct {
	machineId string
	posthog   posthog.Client
	start     time.Time
}

var instance *ClientTelemetry

func Init() {
	instance = ForClient()
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

func BootstrapSuccess(elapsed time.Duration, clusterId string) {
	if instance != nil {
		command, arguments, flags := getCommandAndArgs()
		event := instance.getBaseEvent("bootstrap_success")
		event.Properties["original_command"] = command
		event.Properties["arguments"] = arguments
		event.Properties["flags"] = flags
		event.Properties["execution_time"] = elapsed.Milliseconds()
		event.Properties["cluster_id"] = clusterId
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
	if ph, err := posthog.NewWithConfig(
		apiKey,
		posthog.Config{
			Endpoint: "https://eu.posthog.com",
		},
	); err != nil {
		// TODO ?
		return nil
	} else {
		return &ClientTelemetry{
			machineId: getMachineId(),
			posthog:   ph,
			start:     time.Now(),
		}
	}
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

func (t *ClientTelemetry) getBaseEvent(event string) *posthog.Capture {
	return &posthog.Capture{
		DistinctId: t.machineId,
		Event:      event,
		Properties: map[string]any{
			// event properties:
			"cli_version": config.Version,
			// user properties:
			"$set": map[string]any{
				"cli_version": config.Version,
			},
			"$set_once": map[string]any{
				"type":                "client",
				"initial_cli_version": config.Version,
				"os":                  runtime.GOOS,
				"architecture":        runtime.GOARCH,
			},
		},
	}
}

func (t *ClientTelemetry) report(ctx context.Context, exitCode int, reason string) {
	enabled := false
	clusterId := ""
	operatorVersion := ""
	k8sVersion := ""
	if cfg := client2.ConfigFromContext(ctx); cfg == nil {
		return
	} else if client, e := kubernetes.NewForConfig(cfg); e == nil {
		if namespace, e := client.CoreV1().Namespaces().Get(ctx, "glasskube-system", v1.GetOptions{}); e == nil {
			enabled, clusterId = IsEnabled(namespace.GetAnnotations())
		}
		operatorVersion, _ = clientutils.GetPackageOperatorVersion(ctx)
		if version, e := client.ServerVersion(); e == nil {
			k8sVersion = version.String()
		}
	}
	if !enabled {
		return
	}
	command, arguments, flags := getCommandAndArgs()
	duration := time.Since(t.start).Milliseconds()

	event := t.getBaseEvent(command)
	event.Properties["arguments"] = arguments
	event.Properties["flags"] = flags
	event.Properties["execution_time"] = duration
	event.Properties["exit_code"] = exitCode
	event.Properties["exit_reason"] = reason
	event.Properties["cluster_id"] = clusterId
	event.Properties["operator_version"] = operatorVersion
	event.Properties["cluster_kubernetes_version"] = k8sVersion

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
