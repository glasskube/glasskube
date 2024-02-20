package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/spf13/cobra"
)

var listCmdOptions = struct {
	ListInstalledOnly bool
	ShowDescription   bool
	ShowLatestVersion bool
	More              bool
}{}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List packages",
	Long: "List packages. By default, all available packages of the given repository are shown, " +
		"as well as their installation status in your cluster.\nYou can choose to only show installed packages.",
	PreRun: cliutils.SetupClientContext(true),
	Run: func(cmd *cobra.Command, args []string) {
		if listCmdOptions.More {
			listCmdOptions.ShowLatestVersion = true
			listCmdOptions.ShowDescription = true
		}

		pkgClient := client.FromContext(cmd.Context())
		listOptions := list.DefaultListOptions
		if listCmdOptions.ListInstalledOnly {
			listOptions |= list.OnlyInstalled
		}
		pkgs, err := list.GetPackagesWithStatus(pkgClient, cmd.Context(), listOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred:\n\n%v\n", err)
			os.Exit(1)
			return
		}
		if listCmdOptions.ListInstalledOnly && len(pkgs) == 0 {
			fmt.Println("There are currently no packages installed in your cluster.\n" +
				"Run \"glasskube help install\" to get started.")
		} else {
			printPackageTable(pkgs)
		}
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&listCmdOptions.ListInstalledOnly, "installed", "i", false,
		"list only installed packages")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ShowDescription, "show-description", false,
		"show the package description")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ShowLatestVersion, "show-latest", false,
		"show the latest version of packages if available")
	listCmd.PersistentFlags().BoolVarP(&listCmdOptions.More, "more", "m", false,
		"show additional information about packages (like --show-description --show-latest)")

	listCmd.MarkFlagsMutuallyExclusive("show-description", "more")
	listCmd.MarkFlagsMutuallyExclusive("show-latest", "more")

	RootCmd.AddCommand(listCmd)
}

func printPackageTable(packages []*list.PackageWithStatus) {
	header := []string{"NAME", "STATUS", "VERSION"}
	if listCmdOptions.ShowLatestVersion {
		header = append(header, "LATEST VERSION")
	}
	if listCmdOptions.ShowDescription {
		header = append(header, "DESCRIPTION")
	}
	err := cliutils.PrintPackageTable(os.Stdout,
		packages,
		header,
		func(pkg *list.PackageWithStatus) []string {
			row := []string{pkg.Name, statusString(*pkg), versionString(*pkg)}
			if listCmdOptions.ShowLatestVersion {
				row = append(row, pkg.LatestVersion)
			}
			if listCmdOptions.ShowDescription {
				row = append(row, pkg.ShortDescription)
			}
			return row
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error displaying the package table:\n%v\n(This is a bug)\n", err)
		os.Exit(1)
	}
}

func statusString(pkg list.PackageWithStatus) string {
	if pkg.Status != nil {
		return pkg.Status.Status
	} else {
		return "Not installed"
	}
}

func versionString(pkg list.PackageWithStatus) string {
	if pkg.Package != nil {
		specVersion := pkg.Package.Spec.PackageInfo.Version
		statusVersion := pkg.Package.Status.Version
		repoVersion := pkg.LatestVersion

		if statusVersion != "" {
			versionAddons := []string{}
			if specVersion != "" && statusVersion != specVersion {
				versionAddons = append(versionAddons, fmt.Sprintf("%v desired", specVersion))
			}
			if repoVersion != "" && statusVersion != repoVersion {
				versionAddons = append(versionAddons, fmt.Sprintf("%v available", repoVersion))
			}
			if len(versionAddons) > 0 {
				return fmt.Sprintf("%v (%v)", statusVersion, strings.Join(versionAddons, ", "))
			} else {
				return statusVersion
			}
		} else if specVersion != "" {
			if specVersion != repoVersion {
				return fmt.Sprintf("%v (%v available)", specVersion, repoVersion)
			} else {
				return specVersion
			}
		} else {
			return "n/a"
		}
	} else {
		return ""
	}
}
