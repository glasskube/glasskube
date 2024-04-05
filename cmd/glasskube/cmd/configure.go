package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/manifestvalues/cli"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/manifest"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

var configureCmdOptions = struct {
	values  []string
	keepOld bool
}{
	keepOld: true,
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
	pkgClient := client.FromContext(ctx)
	k8sClient := kubernetes.NewForConfigOrDie(client.ConfigFromContext(ctx))
	valueResolver := manifestvalues.NewResolver(
		clientadapter.NewPackageClientAdapter(pkgClient),
		clientadapter.NewKubernetesClientAdapter(*k8sClient),
	)
	pkgName := args[0]
	var pkg v1alpha1.Package
	var pkgManifest *v1alpha1.PackageManifest

	if err := pkgClient.Packages().Get(ctx, pkgName, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error getting package: %v\n", err)
		os.Exit(1)
	} else if pkgManifest, err = manifest.GetInstalledManifestForPackage(ctx, pkg); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error getting installed manifest: %v\n", err)
		os.Exit(1)
	}

	if configureCmdOptions.keepOld {
		if len(configureCmdOptions.values) > 0 {
			for name, value := range parseValuesFlag(configureCmdOptions.values) {
				if pkg.Spec.Values == nil {
					pkg.Spec.Values = make(map[string]v1alpha1.ValueConfiguration)
				}
				pkg.Spec.Values[name] = value
			}
		} else {
			if values, err := cli.Configure(*pkgManifest, pkg.Spec.Values); err != nil {
				fmt.Fprintf(os.Stderr, "❌ error during configure: %v\n", err)
				os.Exit(1)
			} else {
				pkg.Spec.Values = values
			}
		}
	} else {
		pkg.Spec.Values = parseValuesFlag(configureCmdOptions.values)
	}

	fmt.Fprintln(os.Stderr, bold("Configuration:"))
	printValueConfigurations(os.Stderr, pkg.Spec.Values)
	if _, err := valueResolver.Resolve(ctx, pkg.Spec.Values); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Some values can not be resolved: %v\n", err)
	}

	if !cliutils.YesNoPrompt("Continue?", true) {
		cancel()
	}

	if err := pkgClient.Packages().Update(ctx, &pkg); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error updating package: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Fprintln(os.Stderr, "✅ configuration changed")
	}
}

func parseValuesFlag(str []string) map[string]v1alpha1.ValueConfiguration {
	result := make(map[string]v1alpha1.ValueConfiguration)
	for _, s := range str {
		split := strings.SplitN(s, "=", 2)
		var key, value string
		key = split[0]
		if len(split) > 1 {
			value = split[1]
		}
		result[key] = v1alpha1.ValueConfiguration{Value: &value}
	}
	return result
}

func init() {
	configureCmd.Flags().StringArrayVar(&configureCmdOptions.values, "value", configureCmdOptions.values,
		"set a value via flag (can be used multiple times)")
	configureCmd.Flags().BoolVar(&configureCmdOptions.keepOld, "keep-old", configureCmdOptions.keepOld,
		"set this to false erase any values not specified via --value")
	RootCmd.AddCommand(configureCmd)
}
