package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List packages",
	Long: `List packages. By default, all available packages of the given repository are shown, as well as their installation status in your cluster.
You can choose to only show installed packages.`,
	PreRun: cliutils.SetupClientContext,
	Run: func(cmd *cobra.Command, args []string) {
		pkgClient := client.FromContext(cmd.Context())
		pkgs, err := list.GetPackagesWithStatus(pkgClient, cmd.Context(), config.ListInstalledOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred:\n\n%v\n", err)
			os.Exit(1)
			return
		}
		if config.ListInstalledOnly && len(pkgs) == 0 {
			fmt.Println("There are currently no packages installed in your cluster.\n" +
				"Run \"glasskube help install\" to get started.")
		} else {
			printPackageTable(pkgs, config.Verbose)
		}
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&config.ListInstalledOnly, "installed", "i", false,
		"list only installed packages")
	listCmd.PersistentFlags().BoolVar(&config.Verbose, "show-description", false,
		"show additional information to the packages")
	RootCmd.AddCommand(listCmd)
}

func printPackageTable(packages []*list.PackageTeaserWithStatus, verbose bool) {
	header := make([]string, 3)
	header[0] = "NAME"
	header[1] = "STATUS"
	if verbose {
		header[2] = "DESCRIPTION"
	}
	cliutils.PrintPackageTable(os.Stdout,
		packages,
		header,
		func(pkg *list.PackageTeaserWithStatus) []string {
			row := make([]string, 3)
			row[0] = pkg.PackageName
			var statusStr string
			if pkg.Status == nil {
				statusStr = "Not installed"
			} else {
				statusStr = pkg.Status.Status
			}
			row[1] = statusStr
			if verbose {
				row[2] = pkg.ShortDescription
			}
			return row
		})
}
