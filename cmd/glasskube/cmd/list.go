package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type ListFormat string

const (
	JSON ListFormat = "json"
	YAML ListFormat = "yaml"
)

func (o *ListFormat) String() string {
	return string(*o)
}

func (o *ListFormat) Set(value string) error {
	switch value {
	case string(JSON), string(YAML):
		*o = ListFormat(value)
		return nil
	default:
		return errors.New(`invalid output format, must be "json" or "yaml"`)
	}
}

func (o *ListFormat) Type() string {
	return "string"
}

type ListCmdOptions struct {
	ListInstalledOnly bool
	ListOutdatedOnly  bool
	ShowDescription   bool
	ShowLatestVersion bool
	More              bool
	ListFormat        ListFormat
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
		ctx := cmd.Context()
		if listCmdOptions.More {
			listCmdOptions.ShowLatestVersion = true
			listCmdOptions.ShowDescription = true
		}
		pkgs, err := list.NewLister(ctx).GetPackagesWithStatus(ctx, listCmdOptions.toListOptions())
		if err != nil {
			fmt.Fprintf(os.Stderr, "❗ An error occurred listing packages: %v\n", err)
			if len(pkgs) == 0 {
				cliutils.ExitWithError()
			} else {
				fmt.Fprint(os.Stderr, "⚠️  The table shown below may be incomplete due to the error above.\n\n")
			}
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
			if listCmdOptions.ListFormat == JSON {
				printPackageJSON(pkgs)
			} else if listCmdOptions.ListFormat == YAML {
				printPackageYAML(pkgs)
			} else {
				printPackageTable(pkgs)
			}
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
	listCmd.PersistentFlags().VarP((&listCmdOptions.ListFormat), "output", "o", "output format (json, yaml, etc.)")

	listCmd.MarkFlagsMutuallyExclusive("show-description", "more")
	listCmd.MarkFlagsMutuallyExclusive("show-latest", "more")

	RootCmd.AddCommand(listCmd)
}

func printPackageTable(packages []*list.PackageWithStatus) {
	header := []string{"NAME", "STATUS", "VERSION", "AUTO-UPDATE"}
	if listCmdOptions.ShowLatestVersion {
		header = append(header, "LATEST VERSION")
	}
	header = append(header, "REPOSITORY")
	if listCmdOptions.ShowDescription {
		header = append(header, "DESCRIPTION")
	}
	err := cliutils.PrintTable(os.Stdout,
		packages,
		header,
		func(pkg *list.PackageWithStatus) []string {
			row := []string{pkg.Name, statusString(*pkg), versionString(*pkg), clientutils.AutoUpdateString(pkg.Package, "")}
			if listCmdOptions.ShowLatestVersion {
				row = append(row, pkg.LatestVersion)
			}
			s := make([]string, len(pkg.Repos))
			if pkg.Package != nil {
				for i, r := range pkg.Repos {
					if pkg.Package.Spec.PackageInfo.RepositoryName == r {
						s[i] = fmt.Sprintf("%v (used)", r)
					} else {
						s[i] = r
					}
				}
			} else {
				s = pkg.Repos
			}
			row = append(row, strings.Join(s, ", "))
			if listCmdOptions.ShowDescription {
				row = append(row, pkg.ShortDescription)
			}
			return row
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error displaying the package table:\n%v\n(This is a bug)\n", err)
		cliutils.ExitWithError()
	}
}

func printPackageJSON(packages []*list.PackageWithStatus) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	err := enc.Encode(packages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling data to JSON: %v\n", err)
		cliutils.ExitWithError()
	}
}

func printPackageYAML(packages []*list.PackageWithStatus) {

	for i, pkg := range packages {
		yamlData, err := yaml.Marshal(pkg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error marshaling data to YAML: %v\n", err)
			cliutils.ExitWithError()
		}

		if i > 0 {
			fmt.Println("---")
		}

		fmt.Println(string(yamlData))
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
