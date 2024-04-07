package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/spf13/cobra"
)

type ListCmdOptions struct {
	ListInstalledOnly bool
	ListOutdatedOnly  bool
	ShowDescription   bool
	ShowLatestVersion bool
	More              bool
}

func (o ListCmdOptions) toListOptions() list.ListOptions {
	return list.ListOptions{
		OnlyInstalled: o.ListInstalledOnly,
		OnlyOutdated:  o.ListOutdatedOnly,
	}
}

var listCmdOptions = ListCmdOptions{}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List packages",
	Long: "List packages. By default, all available packages of the given repository are shown, " +
		"as well as their installation status in your cluster.\nYou can choose to only show installed packages.",
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		if listCmdOptions.More {
			listCmdOptions.ShowLatestVersion = true
			listCmdOptions.ShowDescription = true
		}

		pkgClient := client.FromContext(cmd.Context())
		pkgs, err := list.GetPackagesWithStatus(pkgClient, cmd.Context(), listCmdOptions.toListOptions())
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred:\n\n%v\n", err)
			os.Exit(1)
			return
		}
		if len(pkgs) == 0 {
			if listCmdOptions.ListOutdatedOnly {
				fmt.Fprintln(os.Stderr, "All installed packages are up-to-date.")
			} else if listCmdOptions.ListInstalledOnly {
				fmt.Fprintln(os.Stderr, "There are currently no packages installed in your cluster.\n"+
					"Run \"glasskube help install\" to get started.")
			} else {
				fmt.Fprintln(os.Stderr, "No packages found. This is probably a bug.")
			}
		} else {
			printPackageTable(pkgs)
		}
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&listCmdOptions.ListInstalledOnly, "installed", "i", false,
		"list only installed packages")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ListOutdatedOnly, "outdated", false,
		"list only outdated packages")
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
	header := []string{"NAME", "STATUS", "VERSION", "AUTO-UPDATE"}
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
			row := []string{pkg.Name, statusString(*pkg), versionString(*pkg), clientutils.AutoUpdateString(pkg.Package, "")}
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
			if statusVersion != specVersion {
				versionAddons = append(versionAddons, fmt.Sprintf("%v desired", specVersion))
			}
			if repoVersion != "" && semver.IsUpgradable(statusVersion, repoVersion) {
				versionAddons = append(versionAddons, fmt.Sprintf("%v available", repoVersion))
			}
			if len(versionAddons) > 0 {
				return fmt.Sprintf("%v (%v)", statusVersion, strings.Join(versionAddons, ", "))
			} else {
				return statusVersion
			}
		}
		if specVersion != repoVersion {
			return fmt.Sprintf("%v (%v available)", specVersion, repoVersion)
		} else {
			return specVersion
		}
	} else {
		return ""
	}
}
