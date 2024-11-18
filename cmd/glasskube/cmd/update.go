package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/util"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/manifestvalues/cli"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"github.com/glasskube/glasskube/pkg/manifest"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/update"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/cache"
)

var updateCmdOptions = struct {
	cli.ValuesOptions
	Version string
	Yes     bool
	DryRunOptions
	OutputOptions
	NamespaceOptions
	KindOptions
}{
	ValuesOptions: cli.NewOptions(cli.WithKeepOldValuesFlag),
}

var updateCmd = &cobra.Command{
	Use:               "update [<package-name>...]",
	Short:             "Update some or all packages in your cluster",
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: installedPackagesCompletionFunc(&updateCmdOptions.NamespaceOptions, &updateCmdOptions.KindOptions),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		updater := update.NewUpdater(ctx)
		if !rootCmdOptions.NoProgress {
			updater.WithStatusWriter(statuswriter.Spinner())
		}

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

			if pkg, err := getPackageOrClusterPackage(ctx, args[0],
				updateCmdOptions.KindOptions, updateCmdOptions.NamespaceOptions); err != nil {
				fmt.Fprintf(os.Stderr, "Could not get %v: %v\n", args[0], err)
				cliutils.ExitWithError()
			} else {
				tx, err = updater.PrepareForVersion(ctx, pkg, updateCmdOptions.Version)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error in updating the package version : %v\n", err)
					cliutils.ExitWithError()
				}
			}
		} else {
			var updateGetters []update.PackagesGetter
			if len(args) > 0 {
				pkgs := make([]ctrlpkg.Package, len(args))
				for i, name := range args {
					if pkg, err := getPackageOrClusterPackage(ctx, name,
						updateCmdOptions.KindOptions, updateCmdOptions.NamespaceOptions); err != nil {
						fmt.Fprintf(os.Stderr, "Could not get %v: %v\n", name, err)
						cliutils.ExitWithError()
					} else {
						pkgs[i] = pkg
					}
				}
				updateGetters = append(updateGetters, update.GetExact(pkgs))
			} else if updateCmdOptions.Namespace != "" {
				updateGetters = append(updateGetters, update.GetAllPackages(updateCmdOptions.Namespace))
			} else {
				switch updateCmdOptions.Kind {
				case KindClusterPackage:
					updateGetters = append(updateGetters, update.GetAllClusterPackages())
				case KindPackage:
					updateGetters = append(updateGetters, update.GetAllPackages(""))
				default:
					updateGetters = append(updateGetters, update.GetAllClusterPackages(), update.GetAllPackages(""))
				}
			}

			tx, err = updater.Prepare(ctx, updateGetters...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ update preparation failed: %v\n", err)
				cliutils.ExitWithError()
			}
		}

		if tx != nil {
			if len(tx.ConflictItems) > 0 {
				for _, conflictItem := range tx.ConflictItems {
					for _, conflict := range conflictItem.Conflicts {
						fmt.Fprintf(os.Stderr, "❌ Cannot Update %s due to dependency conflicts: %s\n"+
							" (required: %s, actual: %s)\n",
							conflictItem.Package.GetName(), conflict.Actual.Name, conflict.Required.Version, conflict.Actual.Version)
					}
				}
				cliutils.ExitWithError()
			} else if !tx.IsEmpty() {
				printTransaction(*tx)
				if !updateCmdOptions.Yes && !cliutils.YesNoPrompt("Do you want to apply these changes?", false) {
					fmt.Fprintf(os.Stderr, "⛔ Update cancelled. No changes were made.\n")
					cliutils.ExitSuccess()
				}

				for _, item := range tx.Items {
					if item.UpdateRequired() {
						if err := updateConfigurationIfNeeded(ctx, item.Package, item.Version); err != nil {
							fmt.Fprintf(os.Stderr, "❌ error updating configuration for %s: %v\n", item.Package.GetName(), err)
							cliutils.ExitWithError()
						}
					}
				}

				updatedPackages, err := updater.Apply(
					ctx,
					tx,
					update.ApplyUpdateOptions{
						Blocking: true,
						DryRun:   updateCmdOptions.DryRun,
					})
				if err != nil {
					fmt.Fprintf(os.Stderr, "❌ update failed: %v\n", err)
					cliutils.ExitWithError()
				}
				if updateCmdOptions.Output != "" {
					if out, err := clientutils.Format(updateCmdOptions.Output.OutputFormat(),
						updateCmdOptions.ShowAll, updatedPackages...); err != nil {
						fmt.Fprintf(os.Stderr, "❌ failed to marshal output: %v\n", err)
						cliutils.ExitWithError()
					} else {
						fmt.Print(out)
					}
				}
			}
		}

		fmt.Fprintf(os.Stderr, "✅ all packages up-to-date\n")
	},
}

func printTransaction(tx update.UpdateTransaction) {
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', 0)
	if len(tx.Items) > 0 {
		fmt.Fprintf(os.Stderr, "The following packages will be updated:\n")
	}
	for _, item := range tx.Items {
		if item.UpdateRequired() {
			util.Must(fmt.Fprintf(w, " * %v\t%v:\t%v\t-> %v\n",
				item.Package.GetSpec().PackageInfo.Name,
				cache.MetaObjectToName(item.Package),
				item.Package.GetSpec().PackageInfo.Version,
				item.Version,
			))
		} else {
			util.Must(fmt.Fprintf(w, " * %v\t%v:\t%v\t(up-to-date)\n",
				item.Package.GetSpec().PackageInfo.Name,
				cache.MetaObjectToName(item.Package),
				item.Package.GetSpec().PackageInfo.Version,
			))
		}
	}
	for _, req := range tx.Requirements {
		util.Must(fmt.Fprintf(w, " * %v:\t-\t-> %v\n", req.Name, req.Version))
	}
	_ = w.Flush()
	if len(tx.Pruned) > 0 {
		fmt.Fprintf(os.Stderr, "The following packages will be removed:\n")
	}
	for _, req := range tx.Pruned {
		fmt.Fprintf(os.Stderr, " * %v (no longer needed)\n", req.Name)
	}
}

func installedPackagesCompletionFunc(
	nsOpts *NamespaceOptions,
	kindOpts *KindOptions,
) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		dir := cobra.ShellCompDirectiveNoFileComp
		var packages []string

		config, rawConfig, err := kubeconfig.New(config.Kubeconfig)
		if err != nil {
			dir |= cobra.ShellCompDirectiveError
			return nil, dir
		}

		ctx, err := clicontext.SetupContext(cmd.Context(), config, rawConfig)
		if err != nil {
			dir |= cobra.ShellCompDirectiveError
			return nil, dir
		}

		client := cliutils.PackageClient(ctx)

		if (nsOpts == nil || nsOpts.Namespace == "") && (kindOpts == nil || kindOpts.Kind != KindPackage) {
			var list v1alpha1.ClusterPackageList
			if err := client.ClusterPackages().GetAll(ctx, &list); err != nil {
				dir |= cobra.ShellCompDirectiveError
				return nil, dir
			}
			for _, pkg := range list.Items {
				if (toComplete == "" || strings.HasPrefix(pkg.GetName(), toComplete)) && !slices.Contains(args, pkg.GetName()) {
					packages = append(packages, pkg.GetName())
				}
			}
		}

		if kindOpts == nil || kindOpts.Kind != KindClusterPackage {
			ns := ""
			if nsOpts != nil {
				ns = nsOpts.GetActualNamespace(ctx)
			}
			var list v1alpha1.PackageList
			if err := client.Packages(ns).GetAll(ctx, &list); err != nil {
				dir &= cobra.ShellCompDirectiveError
				return nil, dir
			}
			for _, pkg := range list.Items {
				if (toComplete == "" || strings.HasPrefix(pkg.GetName(), toComplete)) && !slices.Contains(args, pkg.GetName()) {
					packages = append(packages, pkg.GetName())
				}
			}
		}

		return packages, dir
	}
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
		dir |= cobra.ShellCompDirectiveError
		return nil, dir
	}
	ctx, err := clicontext.SetupContext(cmd.Context(), config, rawConfig)
	if err != nil {
		dir |= cobra.ShellCompDirectiveError
		return nil, dir
	}
	client := cliutils.PackageClient(ctx)
	repoClient := cliutils.RepositoryClientset(ctx)

	var pkg v1alpha1.ClusterPackage
	if err := client.ClusterPackages().Get(cmd.Context(), packageName, &pkg); err != nil {
		dir |= cobra.ShellCompDirectiveError
		return nil, dir
	}
	var packageIndex repo.PackageIndex
	if err := repoClient.ForPackage(&pkg).FetchPackageIndex(packageName, &packageIndex); err != nil {
		dir |= cobra.ShellCompDirectiveError
		return nil, dir
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

func updateConfigurationIfNeeded(ctx context.Context, pkg ctrlpkg.Package, newVersion string) error {
	newManifest, err := manifest.GetManifestForPackage(ctx, pkg, newVersion)
	if err != nil {
		return fmt.Errorf("error getting manifest for new version: %v", err)
	}

	if updateCmdOptions.ValuesOptions.IsValuesSet() {
		if values, err := updateCmdOptions.ValuesOptions.ParseValues(newManifest, pkg.GetSpec().Values); err != nil {
			return err
		} else {
			pkg.GetSpec().Values = values
		}
	} else if len(newManifest.ValueDefinitions) > 0 || len(pkg.GetSpec().Values) > 0 {
		if cliutils.YesNoPrompt(fmt.Sprintf("Do you want to update the configuration for %s?", pkg.GetName()), false) {
			values, err := cli.Configure(*newManifest,
				cli.WithOldValues(pkg.GetSpec().Values),
				cli.WithUseDefaults(updateCmdOptions.UseDefault),
			)
			if err != nil {
				return fmt.Errorf("error during configuration: %v", err)
			}
			pkg.GetSpec().Values = values
		}
	}

	return nil
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&updateCmdOptions.Version, "version", "v", "",
		"Update to a specific version")
	_ = updateCmd.RegisterFlagCompletionFunc("version", completeUpgradablePackageVersions)
	updateCmd.PersistentFlags().BoolVarP(&updateCmdOptions.Yes, "yes", "y", false,
		"Do not ask for any confirmation")
	updateCmdOptions.OutputOptions.AddFlagsToCommand(updateCmd)
	updateCmdOptions.KindOptions.AddFlagsToCommand(updateCmd)
	updateCmdOptions.NamespaceOptions.AddFlagsToCommand(updateCmd)
	updateCmdOptions.ValuesOptions.AddFlagsToCommand(updateCmd)
	updateCmdOptions.DryRunOptions.AddFlagsToCommand(updateCmd)
	RootCmd.AddCommand(updateCmd)
}
