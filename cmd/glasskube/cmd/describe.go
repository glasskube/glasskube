package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/repo"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:               "describe [package-name]",
	Short:             "Describe a package",
	Long:              "Shows additional information about the given package.",
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		pkg, pkgStatus, manifest, latestVersion, err := describe.DescribePackage(cmd.Context(), pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not describe package %v: %v\n", pkgName, err)
			os.Exit(1)
		}
		bold := color.New(color.Bold).SprintFunc()
		fmt.Printf("%v %v\n", bold("Package Name:"), manifest.Name)
		fmt.Printf("%v %v\n", bold("Short Description:"), stringOrDash(manifest.ShortDescription))
		fmt.Println()
		fmt.Printf("%v \n%v\n", bold("Long Description:"), stringOrDash(manifest.LongDescription))
		if len(manifest.Dependencies) > 0 {
			fmt.Println()
			fmt.Println(bold("Dependencies:"))
			for _, dep := range manifest.Dependencies {
				fmt.Printf(" * %v", dep.Name)
				if len(dep.Version) > 0 {
					fmt.Printf(" (%v)", dep.Version)
				}
				fmt.Println()
			}
		}
		fmt.Println()
		fmt.Printf("%v \n", bold("References:"))
		version := latestVersion
		if pkg != nil {
			version = pkg.Spec.PackageInfo.Version
		}
		if url, err := repo.GetPackageManifestURL("", manifest.Name, version); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not get package manifest url: %v\n", err)
		} else {
			fmt.Printf(" * Glasskube Package Manifest: %v\n", url)
		}
		for _, ref := range manifest.References {
			fmt.Printf(" * %v: %v\n", ref.Label, ref.Url)
		}

		fmt.Printf("\n%v\n", bold("Entrypoints:"))
		if len(manifest.Entrypoints) == 0 {
			fmt.Fprintln(os.Stderr, " * No Entry Points")
		} else {
			for _, i := range manifest.Entrypoints {
				var messageParts []string
				if i.Name != "" {
					messageParts = append(messageParts, fmt.Sprintf("Name: %s", i.Name))
				}
				if i.ServiceName != "" {
					messageParts = append(messageParts, fmt.Sprintf("ServiceName: %s", i.ServiceName))
				}
				if i.Port != 0 {
					messageParts = append(messageParts, fmt.Sprintf("Port: %v", i.Port))
				}
				if i.LocalPort != 0 {
					messageParts = append(messageParts, fmt.Sprintf("LocalPort: %v", i.LocalPort))
				}
				if i.Scheme != "" {
					messageParts = append(messageParts, fmt.Sprintf("Scheme: %s", i.Scheme))
				}
				entrypointMsg := strings.Join(messageParts, ", ")
				fmt.Fprintf(os.Stderr, " * %s\n", entrypointMsg)
			}
		}

		fmt.Printf("\n%v %v\n", bold("Status:"), status(pkgStatus))

	},
}

func stringOrDash(longDesc string) string {
	if len(strings.TrimSpace(longDesc)) > 0 {
		return longDesc
	} else {
		return "–"
	}
}

func status(pkgStatus *client.PackageStatus) string {
	if pkgStatus != nil {
		return pkgStatus.Status
	} else {
		return "Not installed"
	}
}

func init() {
	RootCmd.AddCommand(describeCmd)
}
