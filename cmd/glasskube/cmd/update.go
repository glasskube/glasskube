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
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/update"
	"github.com/spf13/cobra"
)

var updateCmdOptions struct {
	Version string
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

		if len(args) == 1 && updateCmdOptions.Version != "" {
			tx, err := updater.UpdateWithVersion(ctx, args[0], updateCmdOptions.Version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error in updating the package version : %v", err)
				os.Exit(1)
			}

			printTransaction(*tx)

			if !tx.IsEmpty() {
				if !cliutils.YesNoPrompt("Do you want to apply these updates?", false) {
					fmt.Fprintf(os.Stderr, "⛔ Update cancelled. No changes were made.\n")
					os.Exit(0)
				}

				if err := updater.Apply(ctx, tx); err != nil {
					fmt.Fprintf(os.Stderr, "❌ update failed: %v\n", err)
					os.Exit(1)
				}
			}
		} else {
			tx, err := updater.Prepare(ctx, packageNames)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ update preparation failed: %v\n", err)
				os.Exit(1)
			}

			printTransaction(*tx)

			if !tx.IsEmpty() {
				if !cliutils.YesNoPrompt("Do you want to apply these updates?", false) {
					fmt.Fprintf(os.Stderr, "⛔ Update cancelled. No changes were made.\n")
					os.Exit(0)
				}

				if err := updater.Apply(ctx, tx); err != nil {
					fmt.Fprintf(os.Stderr, "❌ update failed: %v\n", err)
					os.Exit(1)
				}
			}
		}

		fmt.Fprintf(os.Stderr, "✅ all packages up-to-date\n")
		os.Exit(0)
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

func init() {
	updateCmd.PersistentFlags().StringVarP(&updateCmdOptions.Version, "version", "v", "",
		"update to a specific version")
	RootCmd.AddCommand(updateCmd)
}
