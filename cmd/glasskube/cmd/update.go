package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/update"
	"github.com/spf13/cobra"
)

var updateCmdOptions struct {
	Version string
	Yes     bool
}

var updateCmd = &cobra.Command{
	Use:               "update [packages...]",
	Short:             "Update some or all packages in your cluster",
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeInstalledPackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client := client.FromContext(ctx)
		packageNames := args
		updater := update.NewUpdater(client).
			WithStatusWriter(statuswriter.Spinner())

		var tx *update.UpdateTransaction
		var err error

		if updateCmdOptions.Version != "" && len(args) > 1 {
			fmt.Fprintf(os.Stderr, "Updating to specific version is only possible for a single package\n")
			cliutils.ExitWithError()
		}
		if len(args) == 1 && updateCmdOptions.Version != "" {
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
			if err := updater.Apply(ctx, tx); err != nil {
				fmt.Fprintf(os.Stderr, "❌ update failed: %v\n", err)
				cliutils.ExitWithError()
			}
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

	var dir cobra.ShellCompDirective
	config, _, err := kubeconfig.New(config.Kubeconfig)
	if err != nil {
		dir &= cobra.ShellCompDirectiveError
		return nil, dir
	}
	client, err := client.New(config)
	if err != nil {
		dir &= cobra.ShellCompDirectiveError
		return nil, dir
	}
	if len(args) == 0 {
		return nil, dir
	}
	packageName := args[0]
	var packageIndex repo.PackageIndex
	if err := repo.FetchPackageIndex("", packageName, &packageIndex); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	if len(args) != 1 {
		return nil, dir
	}
	var pkg v1alpha1.Package
	if err := client.Packages().Get(cmd.Context(), packageName, &pkg); err != nil {
		dir &= cobra.ShellCompDirectiveError
		return nil, dir
	}
	versions := make([]string, 0, len(packageIndex.Versions))
	for _, version := range packageIndex.Versions {
		if toComplete == "" || strings.HasPrefix(version.Version, toComplete) {
			if semver.IsUpgradable(pkg.Spec.PackageInfo.Version, version.Version) {
				versions = append(versions, version.Version)
			}
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
	RootCmd.AddCommand(updateCmd)
}
