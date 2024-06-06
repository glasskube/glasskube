package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/update"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

var updateCmdOptions struct {
	Version string
	Yes     bool
	OutputOptions
}

var updateCmd = &cobra.Command{
	Use:               "update [packages...]",
	Short:             "Update some or all packages in your cluster",
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeInstalledPackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		packageNames := args
		updater := update.NewUpdater(ctx).
			WithStatusWriter(statuswriter.Spinner())

		var tx *update.UpdateTransaction
		var err error

		if updateCmdOptions.Version != "" && len(args) > 1 {
			fmt.Fprintf(os.Stderr, "Updating to specific version is only possible for a single package\n")
			cliutils.ExitWithError()
		}
		if len(args) == 1 && updateCmdOptions.Version != "" {
			if !strings.HasPrefix(updateCmdOptions.Version, "v") {
				updateCmdOptions.Version = "v" + updateCmdOptions.Version
			}
			tx, err = updater.PrepareForVersion(ctx, args[0], updateCmdOptions.Version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error in updating the package version : %v\n", err)
				cliutils.ExitWithError()
			}
		} else {
			tx, err = updater.Prepare(ctx, packageNames)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ update preparation failed: %v\n", err)
				cliutils.ExitWithError()
			}
		}

		if tx != nil && !tx.IsEmpty() {
			printTransaction(*tx)
			if !updateCmdOptions.Yes && !cliutils.YesNoPrompt("Do you want to apply these updates?", false) {
				fmt.Fprintf(os.Stderr, "⛔ Update cancelled. No changes were made.\n")
				cliutils.ExitSuccess()
			}
			updatedPackages, err := updater.Apply(ctx, tx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ update failed: %v\n", err)
				cliutils.ExitWithError()
			}
			handleOutput(updatedPackages)
		}

		fmt.Fprintf(os.Stderr, "✅ all packages up-to-date\n")
	},
}

func printTransaction(tx update.UpdateTransaction) {
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', 0)
	for _, item := range tx.Items {
		if item.UpdateRequired() {
			fmt.Fprintf(w, "%v:\t%v\t-> %v\n",
				item.Package.Name, item.Package.Spec.PackageInfo.Version, item.Version)
		} else {
			fmt.Fprintf(w, "%v:\t%v\t(up-to-date)\n",
				item.Package.Name, item.Package.Spec.PackageInfo.Version)
		}
	}
	for _, req := range tx.Requirements {
		fmt.Fprintf(w, "%v:\t-\t-> %v\n", req.Name, req.Version)
	}
	_ = w.Flush()
}

func handleOutput(pkgs []v1alpha1.Package) {
	if updateCmdOptions.Output == "" {
		return
	}

	var outputData []byte
	var err error
	for i := range pkgs {
		if gvks, _, err := scheme.Scheme.ObjectKinds(&pkgs[i]); err == nil && len(gvks) == 1 {
			pkgs[i].SetGroupVersionKind(gvks[0])
		} else {
			fmt.Fprintf(os.Stderr, "❌ failed to set GVK for package: %v\n", err)
			cliutils.ExitWithError()
		}
	}
	switch updateCmdOptions.Output {
	case OutputFormatJSON:
		outputData, err = json.MarshalIndent(pkgs, "", "  ")
	case OutputFormatYAML:
		var buffer bytes.Buffer
		if len(pkgs) > 1 {
			buffer.Write([]byte("---\n"))
		}
		for i, pkg := range pkgs {
			data, err := yaml.Marshal(pkg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ failed to marshal output: %v\n", err)
				cliutils.ExitWithError()
			}
			buffer.Write(data)
			if i < len(pkgs)-1 {
				buffer.Write([]byte("---\n"))
			}
		}
		outputData = buffer.Bytes()
	default:
		fmt.Fprintf(os.Stderr, "❌ unsupported output format: %v\n", updateCmdOptions.Output)
		cliutils.ExitWithError()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ failed to marshal output: %v\n", err)
		cliutils.ExitWithError()
	}

	fmt.Fprintln(os.Stdout, string(outputData))
}

func completeInstalledPackageNames(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) (packages []string, dir cobra.ShellCompDirective) {
	dir = cobra.ShellCompDirectiveNoFileComp
	config, _, err := kubeconfig.New(config.Kubeconfig)
	if err != nil {
		dir &= cobra.ShellCompDirectiveError
		return
	}
	client, err := client.New(config)
	if err != nil {
		dir &= cobra.ShellCompDirectiveError
		return
	}
	var packageList v1alpha1.PackageList
	if err := client.Packages().GetAll(cmd.Context(), &packageList); err != nil {
		dir &= cobra.ShellCompDirectiveError
		return
	}
	for _, pkg := range packageList.Items {
		if (toComplete == "" || strings.HasPrefix(pkg.GetName(), toComplete)) && !slices.Contains(args, pkg.GetName()) {
			packages = append(packages, pkg.GetName())
		}
	}
	return
}

func completeUpgradablePackageVersions(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	dir := cobra.ShellCompDirectiveNoFileComp
	if len(args) != 1 {
		return nil, dir
	}
	packageName := args[0]

	config, rawConfig, err := kubeconfig.New(config.Kubeconfig)
	if err != nil {
		dir &= cobra.ShellCompDirectiveError
		return nil, dir
	}
	ctx, err := clicontext.SetupContext(cmd.Context(), config, rawConfig)
	if err != nil {
		dir &= cobra.ShellCompDirectiveError
		return nil, dir
	}
	client := cliutils.PackageClient(ctx)
	repoClient := cliutils.RepositoryClientset(ctx)

	var pkg v1alpha1.Package
	if err := client.Packages().Get(cmd.Context(), packageName, &pkg); err != nil {
		dir &= cobra.ShellCompDirectiveError
		return nil, dir
	}
	var packageIndex repo.PackageIndex
	if err := repoClient.ForPackage(pkg).FetchPackageIndex(packageName, &packageIndex); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	versions := make([]string, 0, len(packageIndex.Versions))
	for _, version := range packageIndex.Versions {
		if (toComplete == "" || strings.HasPrefix(version.Version, toComplete)) &&
			semver.IsUpgradable(pkg.Spec.PackageInfo.Version, version.Version) {
			versions = append(versions, version.Version)
		}
	}
	return versions, dir
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&updateCmdOptions.Version, "version", "v", "",
		"update to a specific version")
	_ = updateCmd.RegisterFlagCompletionFunc("version", completeUpgradablePackageVersions)
	updateCmd.PersistentFlags().BoolVarP(&updateCmdOptions.Yes, "yes", "y", false,
		"do not ask for any confirmation")
	updateCmdOptions.OutputOptions.AddFlagsToCommand(updateCmd)
	RootCmd.AddCommand(updateCmd)
}
