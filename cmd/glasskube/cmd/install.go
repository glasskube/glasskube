package cmd

import (
	"fmt"
	"github.com/glasskube/glasskube/api/v1alpha1/condition"
	"github.com/glasskube/glasskube/cmd/glasskube/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [package-name]",
	Short: "Install a package",
	Long:  `Install a package.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgClient, err := client.InitKubeClient(config.Kubeconfig)
		if err != nil {
			return err
		}
		status, err := install.Install(pkgClient, cmd.Context(), args[0])
		if err != nil {
			return err
		}
		if status != nil {
			switch (*status).Status {
			case condition.Ready:
				fmt.Println("Installed successfully.")
			default:
				fmt.Printf("Installation has status %v, reason: %v\nMessage: %v\n",
					(*status).Status, (*status).Reason, (*status).Message)
			}
		} else {
			fmt.Println("Installation status unknown - no error and no status have been observed.")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
