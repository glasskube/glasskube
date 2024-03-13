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
	url           string
	bootstrapType bootstrap.BootstrapType
	latest        bool
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
	PreRun: cliutils.SetupClientContext(false),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)

		if bootstrapCmdOptions.latest {
			updateAvailiable, latestVersion := cliutils.CheckForUpdate()
			if updateAvailiable {
				bootstrapCmdOptions.url = fmt.Sprintf(
					"https://github.com/glasskube/glasskube/releases/download/v%s/manifest-aio.yaml",
					latestVersion,
				)
			} else {
				fmt.Fprintln(os.Stderr, "No updates available, bootstrapping with currently installed version.")
			}
		}

		client := bootstrap.NewBootstrapClient(
			cfg,
			bootstrapCmdOptions.url,
			cmd.Root().Version,
			bootstrapCmdOptions.bootstrapType,
		)
		if err := client.Bootstrap(cmd.Context()); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVarP(&bootstrapCmdOptions.url, "url", "u", "", "URL to fetch the Glasskube operator from")
	bootstrapCmd.Flags().VarP(&bootstrapCmdOptions.bootstrapType, "type", "t", `Type of manifest to use for bootstrapping`)
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.latest, "latest", false, "Fetch and bootstrap the latest version")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "type")

	if config.Version == "dev" {
		bootstrapCmdOptions.latest = true
	}
}
