package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:               "install [package-name]",
	Short:             "Install a package",
	Long:              `Install a package.`,
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.FromContext(cmd.Context())
		packageName := args[0]
		status, err := install.NewInstaller(client).
			WithStatusWriter(statuswriter.Spinner()).
			InstallBlocking(cmd.Context(), packageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
			os.Exit(1)
		}
		if status != nil {
			switch (*status).Status {
			case string(condition.Ready):
				fmt.Printf("✅ %v installed successfully.\n", packageName)
			default:
				fmt.Printf("❌ %v installation has status %v, reason: %v\nMessage: %v\n",
					packageName, (*status).Status, (*status).Reason, (*status).Message)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Installation status unknown - no error and no status have been observed (this is a bug).")
			os.Exit(1)
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
}
