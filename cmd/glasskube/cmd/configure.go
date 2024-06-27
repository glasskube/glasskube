package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/manifestvalues/cli"
	"github.com/glasskube/glasskube/internal/manifestvalues/flags"
	"github.com/glasskube/glasskube/pkg/manifest"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

var configureCmdOptions = struct {
	flags.ValuesOptions
	OutputOptions
	NamespaceOptions
	KindOptions
}{
	ValuesOptions: flags.NewOptions(flags.WithKeepOldValuesFlag),
	KindOptions:   DefaultKindOptions(),
}

var configureCmd = &cobra.Command{
	Use:               "configure <package-name>",
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

	name := args[0]
	pkg, err :=
		getPackageOrClusterPackage(ctx, name, configureCmdOptions.KindOptions, configureCmdOptions.NamespaceOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ could not get resource: %v\n", err)
		cliutils.ExitWithError()
	}

	if configureCmdOptions.IsValuesSet() {
		if values, err := configureCmdOptions.ParseValues(pkg.GetSpec().Values); err != nil {
			fmt.Fprintf(os.Stderr, "❌ invalid values in command line flags: %v\n", err)
			cliutils.ExitWithError()
		} else {
			pkg.GetSpec().Values = values
		}
	} else {
		if pkgManifest, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error getting installed manifest: %v\n", err)
			cliutils.ExitWithError()
		} else if values, err := cli.Configure(*pkgManifest, pkg.GetSpec().Values); err != nil {
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

	if !cliutils.YesNoPrompt("Continue?", true) {
		cancel()
	}

	switch pkg := pkg.(type) {
	case *v1alpha1.ClusterPackage:
		if err := pkgClient.ClusterPackages().Get(ctx, pkg.Name, pkg); err != nil {
			// Don't exit, we can still try to call update ...
			fmt.Fprintf(os.Stderr, "⚠️  error fetching package: %v\n", err)
		}

		if err := pkgClient.ClusterPackages().Update(ctx, pkg); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error updating package: %v\n", err)
			cliutils.ExitWithError()
		}
	case *v1alpha1.Package:
		if err := pkgClient.Packages(pkg.Namespace).Get(ctx, pkg.Name, pkg); err != nil {
			// Don't exit, we can still try to call update ...
			fmt.Fprintf(os.Stderr, "⚠️  error fetching package: %v\n", err)
		}

		if err := pkgClient.Packages(pkg.Namespace).Update(ctx, pkg); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error updating package: %v\n", err)
			cliutils.ExitWithError()
		}
	default:
		fmt.Fprintln(os.Stderr, "❌ invalid state: pkg must be either Package or ClusterPackage")
		cliutils.ExitWithError()
	}

	fmt.Fprintln(os.Stderr, "✅ configuration changed")

	if configureCmdOptions.Output != "" {
		if gvks, _, err := scheme.Scheme.ObjectKinds(pkg); err == nil && len(gvks) == 1 {
			pkg.SetGroupVersionKind(gvks[0])
		}
		var output []byte
		var err error
		switch configureCmdOptions.Output {
		case OutputFormatJSON:
			output, err = json.MarshalIndent(pkg, "", "  ")
		case OutputFormatYAML:
			output, err = yaml.Marshal(pkg)
		default:
			fmt.Fprintf(os.Stderr, "❌ invalid output format: %s\n", configureCmdOptions.Output)
			cliutils.ExitWithError()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ error marshalling output: %v\n", err)
			cliutils.ExitWithError()
		}
		fmt.Println(string(output))
	}
}

func init() {
	configureCmdOptions.ValuesOptions.AddFlagsToCommand(configureCmd)
	configureCmdOptions.OutputOptions.AddFlagsToCommand(configureCmd)
	configureCmdOptions.NamespaceOptions.AddFlagsToCommand(configureCmd)
	configureCmdOptions.KindOptions.AddFlagsToCommand(configureCmd)
	RootCmd.AddCommand(configureCmd)
}
