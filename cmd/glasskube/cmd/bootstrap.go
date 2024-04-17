package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/spf13/cobra"
)

type bootstrapOptions struct {
	url              string
	bootstrapType    bootstrap.BootstrapType
	latest           bool
	disableTelemetry bool
}

var bootstrapCmdOptions = bootstrapOptions{
	bootstrapType: bootstrap.BootstrapTypeAio,
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap Glasskube in a Kubernetes cluster",
	Long: "Bootstraps Glasskube in a Kubernetes cluster, " +
		"thereby installing the Glasskube operator and checking if the installation was successful.",
	Args:   cobra.NoArgs,
	PreRun: cliutils.SetupClientContext(false, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)
		client := bootstrap.NewBootstrapClient(cfg)
		if err := client.Bootstrap(cmd.Context(), bootstrapCmdOptions.asBootstrapOptions()); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			cliutils.ExitWithError()
		}
	},
}

func (o bootstrapOptions) asBootstrapOptions() bootstrap.BootstrapOptions {
	return bootstrap.BootstrapOptions{
		Type:             o.bootstrapType,
		Url:              o.url,
		Latest:           o.latest,
		DisableTelemetry: o.disableTelemetry,
	}
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVarP(&bootstrapCmdOptions.url, "url", "u", "", "URL to fetch the Glasskube operator from")
	bootstrapCmd.Flags().VarP(&bootstrapCmdOptions.bootstrapType, "type", "t", `Type of manifest to use for bootstrapping`)
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.latest, "latest", config.IsDevBuild(),
		"Fetch and bootstrap the latest version")
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.disableTelemetry, "disable-telemetry", false, "Disable telemetry")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "type")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "latest")
}
