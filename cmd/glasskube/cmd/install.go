package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:    "install [package-name]",
	Short:  "Install a package",
	Long:   `Install a package.`,
	Args:   cobra.ExactArgs(1),
	PreRun: cliutils.SetupClientContext,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := client.FromContext(cmd.Context())
		status, err := install.Install(client, cmd.Context(), args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
			os.Exit(1)
		}
		if status != nil {
			switch (*status).Status {
			case string(condition.Ready):
				fmt.Println("Installed successfully.")
			default:
				fmt.Printf("Installation has status %v, reason: %v\nMessage: %v\n",
					(*status).Status, (*status).Reason, (*status).Message)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Installation status unknown - no error and no status have been observed (this is a bug).")
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
