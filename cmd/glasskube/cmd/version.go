package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/constants"
	"github.com/spf13/cobra"
)

var versioncmd = &cobra.Command{
	Use:    "version",
	Short:  "Print the version of glasskube and package-operator",
	Long:   `Print the version of glasskube and package-operator`,
	PreRun: cliutils.SetupClientContext(false, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue(constants.GlasskubeAscii)
		glasskubeVersion := config.Version
		fmt.Fprintf(os.Stderr, "glasskube: v%s\n", glasskubeVersion)
		operatorVersion, err := clientutils.GetPackageOperatorVersion(cmd.Context())
		if err != nil {
			fmt.Fprintf(os.Stderr, "package-operator: not installed\n")
			fmt.Fprintf(os.Stderr, "Glasskube is not yet bootstrapped. Use 'glasskube bootstrap' to get started.\n")
			cliutils.ExitWithError()
		} else {
			fmt.Fprintf(os.Stderr, "package-operator: %s\n", operatorVersion)
		}
	},
}

func init() {
	RootCmd.AddCommand(versioncmd)
}
