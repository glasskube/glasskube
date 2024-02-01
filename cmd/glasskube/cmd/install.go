package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/install"
)

type installOptions struct {
	background bool
}

var installCmdOptions = installOptions{
	background: false,
}

var installCmd = &cobra.Command{
	Use:               "install [package-name]",
	Short:             "Install a package",
	Long:              `Install a package.`,
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.FromContext(cmd.Context())
		status, err := install.InstallPackage(client, cmd.Context(), args[0], installCmdOptions.background)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
			os.Exit(1)
		}

		// TODO: not sure if I should put a return or os.Exit(0) here instead of this
		// heavily nested if/else - jalavosus
		if installCmdOptions.background {
			fmt.Println("Installation running in background, check back for status in a bit.")
		} else {
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
		}
	},
}

func completeAvailablePackageNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	index, err := repo.FetchPackageRepoIndex(cmd.Context(), "")
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}
	names := make([]string, 0, len(index.Packages))
	for _, pkg := range index.Packages {
		if toComplete == "" || strings.HasPrefix(pkg.Name, toComplete) {
			names = append(names, pkg.Name)
		}
	}
	return names, 0
}

func init() {
	RootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&installCmdOptions.background, "background", "b", false, "Run command in background instead of synchronously")
}
