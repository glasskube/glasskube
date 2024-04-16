package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/telemetry/annotations"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
)

var telemetryStatusCmd = &cobra.Command{
	Use:    "status",
	Args:   cobra.NoArgs,
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		rawConfig := client.RawConfigFromContext(ctx)
		bold := color.New(color.Bold).SprintfFunc()
		clientset := kubernetes.NewForConfigOrDie(client.ConfigFromContext(ctx))
		var status string
		if ns, err := clientset.CoreV1().Namespaces().Get(ctx, "glasskube-system", v1.GetOptions{}); err != nil {
			fmt.Fprintf(os.Stderr, "error getting telemetry status: %v\n", err)
			cliutils.ExitWithError()
		} else if annotations.IsTelemetryEnabled(ns.Annotations) {
			status = "enabled"
		} else {
			status = "disabled"
		}
		fmt.Fprintf(os.Stderr, "\nGlasskube telemetry is %v for cluster %v.\n\n"+
			"Run \"glasskube help telemetry\" for more information.\n",
			bold(status), rawConfig.CurrentContext)
	},
}

var telemetryCmd = &cobra.Command{
	Use:       "telemetry <enable|disable>",
	Short:     "View and modify telemetry settings",
	ValidArgs: []string{"enable", "disable"},
	Args:      cobra.ExactArgs(1),
	PreRun:    cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		rawConfig := client.RawConfigFromContext(ctx)
		bold := color.New(color.Bold).SprintfFunc()
		enabled := args[0] == "enable"
		clientset := kubernetes.NewForConfigOrDie(client.ConfigFromContext(ctx))
		if _, err := clientset.CoreV1().Namespaces().Apply(ctx,
			corev1.Namespace("glasskube-system").
				WithAnnotations(map[string]string{annotations.TelemetryEnabledAnnotation: strconv.FormatBool(enabled)}),
			v1.ApplyOptions{FieldManager: "glasskube-telemetry", Force: true}); err != nil {
			fmt.Fprintf(os.Stderr, "error updating telemetry annotations: %v\n", err)
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "\nGlasskube telemetry is now %v for cluster %v.\n",
			bold(args[0]+"d"), rawConfig.CurrentContext)
	},
}

func init() {
	telemetryCmd.AddCommand(telemetryStatusCmd)
	RootCmd.AddCommand(telemetryCmd)
}
