package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/manifestvalues/cli"
	"github.com/glasskube/glasskube/internal/manifestvalues/flags"
	"github.com/glasskube/glasskube/pkg/manifest"
	"github.com/spf13/cobra"
)

var configureCmdOptions = struct{ flags.ValuesOptions }{
	ValuesOptions: flags.NewOptions(flags.WithKeepOldValuesFlag),
}

var configureCmd = &cobra.Command{
	Use:               "configure [package-name]",
	Short:             "Configure a package",
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run:               runConfigure,
	ValidArgsFunction: completeInstalledPackageNames,
}

func runConfigure(cmd *cobra.Command, args []string) {
	bold := color.New(color.Bold).SprintFunc()
	ctx := cmd.Context()
	pkgClient := cliutils.PackageClient(ctx)
	valueResolver := cliutils.ValueResolver(ctx)
	pkgName := args[0]
	var pkg v1alpha1.Package

	if err := pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error getting package: %v\n", err)
		cliutils.ExitWithError()
	}

	if configureCmdOptions.IsValuesSet() {
		if values, err := configureCmdOptions.ParseValues(pkg.Spec.Values); err != nil {
			fmt.Fprintf(os.Stderr, "❌ invalid values in command line flags: %v\n", err)
			cliutils.ExitWithError()
		} else {
			pkg.Spec.Values = values
		}
	} else {
		if pkgManifest, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error getting installed manifest: %v\n", err)
			cliutils.ExitWithError()
		} else if values, err := cli.Configure(*pkgManifest, pkg.Spec.Values); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error during configure: %v\n", err)
			cliutils.ExitWithError()
		} else {
			pkg.Spec.Values = values
		}
	}

	fmt.Fprintln(os.Stderr, bold("Configuration:"))
	printValueConfigurations(os.Stderr, pkg.Spec.Values)
	if _, err := valueResolver.Resolve(ctx, pkg.Spec.Values); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Some values can not be resolved: %v\n", err)
	}

	if !cliutils.YesNoPrompt("Continue?", true) {
		cancel()
	}

	if err := pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		// Don't exit, we can still try to call update ...
		fmt.Fprintf(os.Stderr, "⚠️  error fetching package: %v\n", err)
	}

	if err := pkgClient.Packages().Update(ctx, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error updating package: %v\n", err)
		cliutils.ExitWithError()
	} else {
		fmt.Fprintln(os.Stderr, "✅ configuration changed")
	}
}

func init() {
	configureCmdOptions.ValuesOptions.AddFlagsToCommand(configureCmd)
	RootCmd.AddCommand(configureCmd)
}
