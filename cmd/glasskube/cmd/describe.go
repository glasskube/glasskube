package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/repo"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:               "describe [package-name]",
	Short:             "Describe a package",
	Long:              "Shows additional information about the given package.",
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		pkg, pkgStatus, manifest, latestVersion, err := describe.DescribePackage(cmd.Context(), pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not describe package %v: %v\n", pkgName, err)
			os.Exit(1)
		}
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println(bold("Package:"), nameAndDescription(manifest))
		fmt.Println(bold("Version:"), version(pkg, latestVersion))
		fmt.Println(bold("Status: "), status(pkgStatus))
		if pkg != nil {
			fmt.Println(bold("Auto-Update:"), clientutils.AutoUpdateString(pkg, "Disabled"))
		}

		if len(manifest.Entrypoints) > 0 {
			fmt.Println()
			fmt.Println(bold("Entrypoints:"))
			printEntrypoints(manifest)
		}

		if len(manifest.Dependencies) > 0 {
			fmt.Println()
			fmt.Println(bold("Dependencies:"))
			printDependencies(manifest)
		}

		fmt.Println()
		fmt.Printf("%v \n", bold("References:"))
		printReferences(pkg, manifest, latestVersion)

		trimmedDescription := strings.TrimSpace(manifest.LongDescription)
		if len(trimmedDescription) > 0 {
			fmt.Println()
			fmt.Println(bold("Long Description:"))
			fmt.Println(trimmedDescription)
		}

		if pkg != nil && len(pkg.Spec.Values) > 0 {
			fmt.Println()
			fmt.Println(bold("Configuration:"))
			printValueConfigurations(os.Stdout, pkg.Spec.Values)
		}
	},
}

func printEntrypoints(manifest *v1alpha1.PackageManifest) {
	for _, i := range manifest.Entrypoints {
		var messageParts []string
		if i.Name != "" {
			messageParts = append(messageParts, fmt.Sprint("Name: ", i.Name))
		}
		if i.ServiceName != "" {
			messageParts = append(messageParts, fmt.Sprintf("Remote: %s:%v", i.ServiceName, i.Port))
		}
		var localUrl string
		if i.Scheme != "" {
			localUrl += i.Scheme + "://localhost:"
		} else {
			localUrl += "http://localhost:"
		}
		if i.LocalPort != 0 {
			localUrl += fmt.Sprint(i.LocalPort)
		} else {
			localUrl += fmt.Sprint(i.Port)
		}
		localUrl += "/"
		messageParts = append(messageParts, fmt.Sprint("Local: ", localUrl))
		entrypointMsg := strings.Join(messageParts, ", ")
		fmt.Fprintf(os.Stderr, " * %s\n", entrypointMsg)
	}
}

func printDependencies(manifest *v1alpha1.PackageManifest) {
	for _, dep := range manifest.Dependencies {
		fmt.Printf(" * %v", dep.Name)
		if len(dep.Version) > 0 {
			fmt.Printf(" (%v)", dep.Version)
		}
		fmt.Println()
	}
}

func printReferences(pkg *v1alpha1.Package, manifest *v1alpha1.PackageManifest, latestVersion string) {
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
}

func printValueConfigurations(w io.Writer, values map[string]v1alpha1.ValueConfiguration) {
	for name, value := range values {
		fmt.Fprintf(w, " * %v: %v\n", name, manifestvalues.ValueAsString(value))
	}
}

func status(pkgStatus *client.PackageStatus) string {
	if pkgStatus != nil {
		switch pkgStatus.Status {
		case string(condition.Ready):
			return color.GreenString(pkgStatus.Status)
		case string(condition.Failed):
			return color.RedString(pkgStatus.Status)
		default:
			return pkgStatus.Status
		}
	} else {
		return color.New(color.Faint).Sprint("Not installed")
	}
}

func version(pkg *v1alpha1.Package, latestVersion string) string {
	if len(latestVersion) == 0 {
		if pkg.Spec.PackageInfo.Version != pkg.Status.Version {
			return fmt.Sprintf("%v (desired: %v)", pkg.Status.Version, pkg.Spec.PackageInfo.Version)
		} else {
			return pkg.Status.Version
		}
	} else {
		return latestVersion
	}
}

func nameAndDescription(manifest *v1alpha1.PackageManifest) string {
	var bld strings.Builder
	_, _ = bld.WriteString(manifest.Name)
	trimmedDescription := strings.TrimSpace(manifest.ShortDescription)
	if len(trimmedDescription) > 0 {
		_, _ = bld.WriteString(" — ")
		_, _ = bld.WriteString(trimmedDescription)
	}
	return bld.String()
}

func init() {
	RootCmd.AddCommand(describeCmd)
}
