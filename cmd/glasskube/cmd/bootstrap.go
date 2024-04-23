package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/spf13/cobra"
)

type bootstrapOptions struct {
	url              string
	bootstrapType    bootstrap.BootstrapType
	latest           bool
	disableTelemetry bool
	force            bool
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
	PreRun: cliutils.SetupClientContext(false, util.Pointer(true)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)
		client := bootstrap.NewBootstrapClient(cfg)

		installedVersion, err := clientutils.GetPackageOperatorVersion(cmd.Context())
		if err != nil {
			IsBootstrapped, err := bootstrap.IsBootstrapped(cmd.Context(), cfg)
			if err != nil && !IsBootstrapped {
				fmt.Printf("error : %v\n", err)
				cliutils.ExitWithError()
			}
		}

		var desiredVersion string
		if bootstrapCmdOptions.url == "" {
			desiredVersion = config.Version
		} else {
			desiredVersion = ""
		}
		if !semver.IsUpgradable(installedVersion, desiredVersion) &&
			installedVersion != "" &&
			installedVersion[1:] != desiredVersion {
			if !cliutils.YesNoPrompt(fmt.Sprintf("Glasskube is already installed in this cluster "+
				"in the newer version %v. You are about to install version %v. This could lead "+
				"to a broken cluster!\nAre you sure that you want to downgrade glasskube "+
				"in this cluster?", installedVersion, desiredVersion), false) {
				fmt.Println("Operation stopped")
				cliutils.ExitWithError()
			}

		}

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
		Force:            o.force,
	}
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVarP(&bootstrapCmdOptions.url, "url", "u", "", "URL to fetch the Glasskube operator from")
	bootstrapCmd.Flags().VarP(&bootstrapCmdOptions.bootstrapType, "type", "t", `Type of manifest to use for bootstrapping`)
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.latest, "latest", config.IsDevBuild(),
		"Fetch and bootstrap the latest version")
	bootstrapCmd.Flags().BoolVarP(&bootstrapCmdOptions.force, "force", "f", bootstrapCmdOptions.force,
		"Do not bail out if pre-checks fail")
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.disableTelemetry, "disable-telemetry", false, "Disable telemetry")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "type")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "latest")
}
