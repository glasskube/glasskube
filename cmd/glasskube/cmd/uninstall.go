package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/uninstall"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:    "uninstall [package-name]",
	Short:  "Uninstall a package",
	Long:   `Uninstall a package.`,
	Args:   cobra.ExactArgs(1),
	PreRun: cliutils.SetupClientContext,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.FromContext(cmd.Context())
		ok, err := uninstall.Uninstall(client, cmd.Context(), args[0], config.ForceUninstall)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred during uninstallation:\n\n%v\n", err)
			os.Exit(1)
			return
		}
		if ok {
			fmt.Println("Uninstalled successfully.")
		} else {
			fmt.Println("Uninstallation cancelled.")
		}
	},
}

func init() {
	uninstallCmd.PersistentFlags().BoolVar(&config.ForceUninstall, "force", false,
		"skip the confirmation question and uninstall right away")
	RootCmd.AddCommand(uninstallCmd)
}
