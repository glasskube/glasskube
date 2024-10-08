package cmd

import (
	"fmt"
	"maps"
	"os"

	"github.com/glasskube/glasskube/internal/clientutils"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/manifestvalues/cli"
	"github.com/glasskube/glasskube/pkg/manifest"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var configureCmdOptions = struct {
	cli.ValuesOptions
	OutputOptions
	NamespaceOptions
	KindOptions
	DryRunOptions
}{
	ValuesOptions: cli.NewOptions(cli.WithKeepOldValuesFlag),
	KindOptions:   DefaultKindOptions(),
}

var configureCmd = &cobra.Command{
	Use:    "configure <package-name>",
	Short:  "Configure a package",
	Args:   cobra.ExactArgs(1),
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run:    runConfigure,
	ValidArgsFunction: installedPackagesCompletionFunc(
		&configureCmdOptions.NamespaceOptions,
		&configureCmdOptions.KindOptions,
	),
}

func runConfigure(cmd *cobra.Command, args []string) {
	bold := color.New(color.Bold).SprintFunc()
	ctx := cmd.Context()
	pkgClient := cliutils.PackageClient(ctx)
	valueResolver := cliutils.ValueResolver(ctx)
	name := args[0]

	opts := metav1.UpdateOptions{}
	if configureCmdOptions.DryRun {
		opts.DryRun = []string{metav1.DryRunAll}
		fmt.Fprintln(os.Stderr,
			"🔎 Dry-run mode is enabled. Nothing will be changed.")
	}

	pkg, err :=
		getPackageOrClusterPackage(ctx, name, configureCmdOptions.KindOptions, configureCmdOptions.NamespaceOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ could not get resource: %v\n", err)
		cliutils.ExitWithError()
	}

	pkgManifest, err := manifest.GetInstalledManifestForPackage(ctx, pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ error getting installed manifest: %v\n", err)
		cliutils.ExitWithError()
	}

	if configureCmdOptions.IsValuesSet() {
		if values, err := configureCmdOptions.ParseValues(pkgManifest, pkg.GetSpec().Values); err != nil {
			fmt.Fprintf(os.Stderr, "❌ invalid values in command line flags: %v\n", err)
			cliutils.ExitWithError()
		} else {
			pkg.GetSpec().Values = values
		}
	} else {
		if len(pkgManifest.ValueDefinitions) == 0 {
			fmt.Fprintln(os.Stderr, "❌ this package has no configuration values")
			cliutils.ExitWithError()
		} else if values, err := cli.Configure(*pkgManifest,
			cli.WithOldValues(pkg.GetSpec().Values),
			cli.WithUseDefaults(configureCmdOptions.UseDefault),
		); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error during configure: %v\n", err)
			cliutils.ExitWithError()
		} else {
			pkg.GetSpec().Values = values
		}
	}

	fmt.Fprintln(os.Stderr, bold("Configuration:"))
	printValueConfigurations(os.Stderr, pkg.GetSpec().Values)
	if _, err := valueResolver.Resolve(ctx, pkg.GetSpec().Values); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Some values can not be resolved: %v\n", err)
	}

	if !configureCmdOptions.DryRun {
		if !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}
	}

	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		values := maps.Clone(pkg.Spec.Values)
		if err := pkgClient.ClusterPackages().Get(ctx, pkg.Name, pkg); err != nil {
			// Don't exit, we can still try to call update ...
			fmt.Fprintf(os.Stderr, "⚠️  error fetching package: %v\n", err)
		}
		pkg.Spec.Values = values

		if err := pkgClient.ClusterPackages().Update(ctx, pkg, opts); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error updating package: %v\n", err)
			cliutils.ExitWithError()
		}
	case *v1alpha1.Package:
		values := maps.Clone(pkg.Spec.Values)
		if err := pkgClient.Packages(pkg.Namespace).Get(ctx, pkg.Name, pkg); err != nil {
			// Don't exit, we can still try to call update ...
			fmt.Fprintf(os.Stderr, "⚠️  error fetching package: %v\n", err)
		}
		pkg.Spec.Values = values

		if err := pkgClient.Packages(pkg.Namespace).Update(ctx, pkg, opts); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error updating package: %v\n", err)
			cliutils.ExitWithError()
		}
	default:
		fmt.Fprintln(os.Stderr, "❌ invalid state: pkg must be either Package or ClusterPackage")
		cliutils.ExitWithError()
	}

	if configureCmdOptions.DryRun {
		fmt.Fprintln(os.Stderr, "✅ valid configuration but nothing has been changed")
	} else {
		fmt.Fprintln(os.Stderr, "✅ configuration changed")
	}

	if configureCmdOptions.Output != "" {
		if out, err := clientutils.Format(configureCmdOptions.Output.OutputFormat(),
			configureCmdOptions.ShowAll, pkg); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error marshalling output: %v\n", err)
			cliutils.ExitWithError()
		} else {
			fmt.Print(out)
		}
	}
}

func init() {
	configureCmdOptions.ValuesOptions.AddFlagsToCommand(configureCmd)
	configureCmdOptions.OutputOptions.AddFlagsToCommand(configureCmd)
	configureCmdOptions.NamespaceOptions.AddFlagsToCommand(configureCmd)
	configureCmdOptions.KindOptions.AddFlagsToCommand(configureCmd)
	configureCmdOptions.DryRunOptions.AddFlagsToCommand(configureCmd)
	RootCmd.AddCommand(configureCmd)
}
