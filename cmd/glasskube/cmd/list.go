package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type ListCmdOptions struct {
	ListInstalledOnly bool
	ListOutdatedOnly  bool
	ShowDescription   bool
	ShowLatestVersion bool
	More              bool
	OutputOptions
	KindOptions
}

func (o ListCmdOptions) toListOptions() list.ListOptions {
	return list.ListOptions{
		OnlyInstalled: o.ListInstalledOnly,
		OnlyOutdated:  o.ListOutdatedOnly,
	}
}

var listCmdOptions = ListCmdOptions{
	KindOptions: DefaultKindOptions(),
}

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
		lister := list.NewListerWithRepoCache(ctx)
		var clPkgs []*list.PackageWithStatus
		var pkgs []*list.PackagesWithStatus
		var err error
		if listCmdOptions.Kind != KindPackage {
			clPkgs, err = lister.GetClusterPackagesWithStatus(ctx, listCmdOptions.toListOptions())
			handleListErr(len(clPkgs), err, "clusterpackages")
		}
		if listCmdOptions.Kind != KindClusterPackage {
			pkgs, err = lister.GetPackagesWithStatus(ctx, listCmdOptions.toListOptions())
			handleListErr(len(pkgs), err, "packages")
		}
		noPkgs := len(pkgs) == 0 && listCmdOptions.Kind != KindClusterPackage
		noClPkgs := len(clPkgs) == 0 && listCmdOptions.Kind != KindPackage
		if noPkgs {
			handleEmptyList("packages")
		} else if len(pkgs) > 0 {
			printPackageTable(pkgs)
			if listCmdOptions.Kind != KindPackage {
				fmt.Fprintln(os.Stderr, "")
			}
		}
		if noClPkgs {
			handleEmptyList("clusterpackages")
		} else if len(clPkgs) > 0 {
			if listCmdOptions.Output == OutputFormatJSON {
				printPackageJSON(clPkgs)
			} else if listCmdOptions.Output == OutputFormatYAML {
				printPackageYAML(clPkgs)
			} else {
				printClusterPackageTable(clPkgs)
			}
		}
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&listCmdOptions.ListInstalledOnly, "installed", "i", false,
		"list only installed (cluster-)packages")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ListOutdatedOnly, "outdated", false,
		"list only outdated (cluster-)packages")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ShowDescription, "show-description", false,
		"show the (cluster-)package description")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ShowLatestVersion, "show-latest", false,
		"show the latest version of (cluster-)packages if available")
	listCmd.PersistentFlags().BoolVarP(&listCmdOptions.More, "more", "m", false,
		"show additional information about (cluster-)packages (like --show-description --show-latest)")
	listCmdOptions.OutputOptions.AddFlagsToCommand(listCmd)
	listCmdOptions.KindOptions.AddFlagsToCommand(listCmd)

	listCmd.MarkFlagsMutuallyExclusive("show-description", "more")
	listCmd.MarkFlagsMutuallyExclusive("show-latest", "more")

	RootCmd.AddCommand(listCmd)
}

func handleListErr(listLen int, err error, resource string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "❗ An error occurred listing %s: %v\n", resource, err)
		if listLen == 0 {
			cliutils.ExitWithError()
		} else {
			fmt.Fprint(os.Stderr, "⚠️  The table shown below may be incomplete due to the error above.\n\n")
		}
	}
}

func handleEmptyList(resource string) {
	if listCmdOptions.ListOutdatedOnly {
		fmt.Fprintf(os.Stderr, "All installed %s are up-to-date.\n", resource)
	} else if listCmdOptions.ListInstalledOnly {
		fmt.Fprintf(os.Stderr, "There are currently no %s installed in your cluster.\n"+
			"Run \"glasskube help install\" to get started.\n", resource)
	} else {
		fmt.Fprintf(os.Stderr, "No %s found in the available repositories.\n", resource)
	}
}

func printClusterPackageTable(packages []*list.PackageWithStatus) {
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
			row := []string{pkg.Name, statusString(*pkg), versionString(*pkg),
				clientutils.AutoUpdateString(pkg.ClusterPackage, "")}
			if listCmdOptions.ShowLatestVersion {
				row = append(row, pkg.LatestVersion)
			}
			s := make([]string, len(pkg.Repos))
			if pkg.ClusterPackage != nil {
				for i, r := range pkg.Repos {
					if pkg.ClusterPackage.Spec.PackageInfo.RepositoryName == r {
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
		fmt.Fprintf(os.Stderr, "There was an error displaying the clusterpackage table:\n%v\n(This is a bug)\n", err)
		cliutils.ExitWithError()
	}
}

func printPackageTable(packages []*list.PackagesWithStatus) {
	header := []string{"PACKAGENAME", "NAMESPACE", "NAME", "STATUS", "VERSION", "AUTO-UPDATE"}
	if listCmdOptions.ShowLatestVersion {
		header = append(header, "LATEST VERSION")
	}
	header = append(header, "REPOSITORY")
	if listCmdOptions.ShowDescription {
		header = append(header, "DESCRIPTION")
	}

	var flattenedPkgs []*list.PackageWithStatus
	for _, pkgs := range packages {
		if len(pkgs.Packages) == 0 {
			flattenedPkgs = append(flattenedPkgs, &list.PackageWithStatus{
				MetaIndexItem: pkgs.MetaIndexItem,
			})
		} else {
			flattenedPkgs = append(flattenedPkgs, pkgs.Packages...)
		}
	}

	err := cliutils.PrintTable(os.Stdout,
		flattenedPkgs,
		header,
		func(pkg *list.PackageWithStatus) []string {
			row := []string{pkg.Name, pkgNamespaceString(*pkg), pkgNameString(*pkg), statusString(*pkg), versionString(*pkg),
				clientutils.AutoUpdateString(pkg.Package, "")}
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

func pkgNamespaceString(pkg list.PackageWithStatus) string {
	if pkg.Package != nil {
		return pkg.Package.Namespace
	} else {
		return ""
	}
}

func pkgNameString(pkg list.PackageWithStatus) string {
	if pkg.Package != nil {
		return pkg.Package.Name
	} else {
		return ""
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
	var p ctrlpkg.Package
	if pkg.ClusterPackage != nil {
		p = pkg.ClusterPackage
	} else if pkg.Package != nil {
		p = pkg.Package
	}
	if pkg.ClusterPackage != nil || pkg.Package != nil {
		specVersion := p.GetSpec().PackageInfo.Version
		statusVersion := p.GetStatus().Version
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
