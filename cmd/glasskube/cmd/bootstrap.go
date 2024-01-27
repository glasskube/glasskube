package cmd

import (
	"fmt"
	"github.com/glasskube/glasskube/cmd/glasskube/config"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/spf13/cobra"
	"os"
)

type bootstrapOptions struct {
	url string
}

var bootstrapCmdOptions bootstrapOptions

var bootstrapCmd = &cobra.Command{
	Use:    "bootstrap",
	Short:  "Bootstrap Glasskube in a Kubernetes cluster",
	Long:   `Bootstraps Glasskube in a Kubernetes cluster, thereby installing the Glasskube operator and checking if the installation was successful.`,
	Args:   cobra.ExactArgs(0),
	PreRun: cliutils.SetupClientContext,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := bootstrap.NewBootstrapClient(cmd.Root().Version, config.Kubeconfig, bootstrapCmdOptions.url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred during client initialization:\n\n%v\n", err)
			os.Exit(1)
		}

		err = client.Bootstrap()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVarP(&bootstrapCmdOptions.url, "url", "u", "", "URL to fetch the Glasskube operator from")
}
