package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/clientutils"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/spf13/cobra"
)

var versioncmd = &cobra.Command{
	Use:    "version",
	Short:  "Print the version of glasskube and package-operator",
	Long:   `Print the version of glasskube and package-operator`,
	PreRun: cliutils.SetupClientContext(false, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		glasskubeVersion := config.Version
		fmt.Fprintf(os.Stderr, "glasskube: v%s\n", glasskubeVersion)
		operatorVersion, err := clientutils.GetPackageOperatorVersion(cmd.Context())
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ no deployments found in the glasskube-system namespace\n")
			cliutils.ExitWithError()
		} else {
			fmt.Fprintf(os.Stderr, "package-operator: %s\n", operatorVersion)
		}
	},
}

func init() {
	RootCmd.AddCommand(versioncmd)
}
