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
	ShowMessage       bool
	More              bool
	Repository        string
	packageName       string
	OutputOptions
	KindOptions
	NamespaceOptions
}

func (o ListCmdOptions) toListOptions() list.ListOptions {
	return list.ListOptions{
		OnlyInstalled: o.ListInstalledOnly,
		OnlyOutdated:  o.ListOutdatedOnly,
		Repository:    o.Repository,
		PackageName:   o.packageName,
		Namespace:     o.Namespace,
	}
}

var listCmdOptions = ListCmdOptions{
	KindOptions: DefaultKindOptions(),
}

var listCmd = &cobra.Command{
	Use:     "list [<package-name>]",
	Aliases: []string{"ls", "l"},
	Short:   "List packages",
	Long: "List packages. By default, all available packages of the given repository are shown, " +
		"as well as their installation status in your cluster.\nYou can choose to only show installed packages.",
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Args:   cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if listCmdOptions.More {
			listCmdOptions.ShowLatestVersion = true
			listCmdOptions.ShowDescription = true
			listCmdOptions.ShowMessage = true
		}
		if len(args) > 0 {
			listCmdOptions.packageName = args[0]
		}
		if listCmdOptions.Kind == KindClusterPackage &&
			(listCmdOptions.packageName != "" || listCmdOptions.Namespace != "") {
			fmt.Fprintf(os.Stderr, "Argument [<package-name>] or flag [--namespace] not supported with kind %s.\n",
				KindClusterPackage)
			cliutils.ExitWithError()
		}
		lister := list.NewListerWithRepoCache(ctx)
		var clPkgs []*list.PackageWithStatus
		var pkgs []*list.PackagesWithStatus
		var err error
		if listCmdOptions.Kind != KindPackage && listCmdOptions.packageName == "" && listCmdOptions.Namespace == "" {
			clPkgs, err = lister.GetClusterPackagesWithStatus(ctx, listCmdOptions.toListOptions())
			handleListErr(len(clPkgs), err, "clusterpackages")
		}
		if listCmdOptions.Kind != KindClusterPackage {
			pkgs, err = lister.GetPackagesWithStatus(ctx, listCmdOptions.toListOptions())
			handleListErr(len(pkgs), err, "packages")
		}
		noPkgs := len(pkgs) == 0 && listCmdOptions.Kind != KindClusterPackage
		noClPkgs := len(clPkgs) == 0 && listCmdOptions.Kind != KindPackage &&
			listCmdOptions.packageName == "" && listCmdOptions.Namespace == ""
		if listCmdOptions.Output == outputFormatJSON {
			printPackageJSON(allPkgs(clPkgs, pkgs))
		} else if listCmdOptions.Output == outputFormatYAML {
			printPackageYAML(allPkgs(clPkgs, pkgs))
		} else {
			if noPkgs {
				handleEmptyList("packages")
			} else if len(pkgs) > 0 {
				printPackageTable(pkgs)
			}
			if noClPkgs {
				handleEmptyList("clusterpackages")
			} else if len(clPkgs) > 0 {
				if len(pkgs) > 0 {
					fmt.Fprintln(os.Stderr, "")
				}
				printClusterPackageTable(clPkgs)
			}
		}
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&listCmdOptions.ListInstalledOnly, "installed", "i", false,
		"List only installed (cluster-)packages")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ListOutdatedOnly, "outdated", false,
		"List only outdated (cluster-)packages")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ShowDescription, "show-description", false,
		"Show the (cluster-)package description")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ShowLatestVersion, "show-latest", false,
		"Show the latest version of (cluster-)packages if available")
	listCmd.PersistentFlags().BoolVar(&listCmdOptions.ShowMessage, "show-message", false,
		"Show the messages of (cluster-)packages")
	listCmd.PersistentFlags().BoolVarP(&listCmdOptions.More, "more", "m", false,
		"Show additional information about (cluster-)packages (like --show-description --show-latest)")
	listCmd.PersistentFlags().StringVarP(&listCmdOptions.Repository, "repository", "r", "",
		"Filter based on the repository provided")
	listCmdOptions.OutputOptions.AddFlagsToCommand(listCmd)
	listCmdOptions.KindOptions.AddFlagsToCommand(listCmd)
	listCmdOptions.NamespaceOptions.AddFlagsToCommand(listCmd)

	listCmd.MarkFlagsMutuallyExclusive("show-description", "more")
	listCmd.MarkFlagsMutuallyExclusive("show-latest", "more")
	listCmd.MarkFlagsMutuallyExclusive("show-message", "more")

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
	fmt.Fprintf(os.Stderr, "No %s found.\n", resource)
}

func allPkgs(clpkgs []*list.PackageWithStatus, pkgs []*list.PackagesWithStatus) []*list.PackageWithStatus {
	result := make([]*list.PackageWithStatus, 0, len(clpkgs)+len(pkgs))
	result = append(result, clpkgs...)
	for _, pkg := range pkgs {
		if len(pkg.Packages) > 0 {
			result = append(result, pkg.Packages...)
		} else {
			result = append(result, &list.PackageWithStatus{
				MetaIndexItem: pkg.MetaIndexItem,
			})
		}
	}
	return result
}

func printClusterPackageTable(packages []*list.PackageWithStatus) {
	header := []string{"NAME", "VERSION", "AUTO-UPDATE", "SUSPENDED"}
	if listCmdOptions.ShowLatestVersion {
		header = append(header, "LATEST VERSION")
	}
	header = append(header, "REPOSITORY")
	if listCmdOptions.ShowDescription {
		header = append(header, "DESCRIPTION")
	}
	header = append(header, "STATUS")
	if listCmdOptions.ShowMessage {
		header = append(header, "MESSAGE")
	}

	err := cliutils.PrintTable(os.Stdout,
		packages,
		header,
		func(pkg *list.PackageWithStatus) []string {
			row := []string{pkg.Name, versionString(*pkg),
				clientutils.AutoUpdateString(pkg.ClusterPackage, "")}
			if pkg.ClusterPackage != nil {
				row = append(row, boolYesNo(pkg.ClusterPackage.Spec.Suspend))
			} else {
				row = append(row, "")
			}
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
			row = append(row, statusString(*pkg))
			if listCmdOptions.ShowMessage {
				row = append(row, messageString(*pkg))
			}
			return row
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error displaying the clusterpackage table:\n%v\n(This is a bug)\n", err)
		cliutils.ExitWithError()
	}
}

func printPackageTable(packages []*list.PackagesWithStatus) {
	header := []string{"PACKAGENAME", "NAMESPACE", "NAME", "VERSION", "AUTO-UPDATE", "SUSPENDED"}
	if listCmdOptions.ShowLatestVersion {
		header = append(header, "LATEST VERSION")
	}
	header = append(header, "REPOSITORY")
	if listCmdOptions.ShowDescription {
		header = append(header, "DESCRIPTION")
	}
	header = append(header, "STATUS")
	if listCmdOptions.ShowMessage {
		header = append(header, "MESSAGE")
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
			row := []string{pkg.Name, pkgNamespaceString(*pkg), pkgNameString(*pkg), versionString(*pkg),
				clientutils.AutoUpdateString(pkg.Package, "")}
			if pkg.Package != nil {
				row = append(row, boolYesNo(pkg.Package.Spec.Suspend))
			} else {
				row = append(row, "")
			}
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
			row = append(row, statusString(*pkg))
			if listCmdOptions.ShowMessage {
				row = append(row, messageString(*pkg))
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

func messageString(pkg list.PackageWithStatus) string {
	if pkg.Status != nil {
		return pkg.Status.Message
	} else {
		return ""
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
