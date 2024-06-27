package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/update"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

var autoUpdateEnabledDisabledOptions = struct {
	Yes, All bool
	KindOptions
	NamespaceOptions
}{
	KindOptions: DefaultKindOptions(),
}

var autoUpdateEnableCmd = &cobra.Command{
	Use:               "enable [...package]",
	Short:             "Enable automatic updates for packages:",
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeInstalledPackageNames,
	Run: runAutoUpdateEnableOrDisable(true,
		"Enable automatic updates for the following packages", "Automatic updates enabled"),
}

var autoUpdateDisableCmd = &cobra.Command{
	Use:               "disable [...package]",
	Short:             "Disable automatic updates for packages:",
	PreRun:            cliutils.SetupClientContext(false, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeInstalledPackageNames,
	Run: runAutoUpdateEnableOrDisable(false,
		"Enable automatic updates for the following packages", "Automatic updates disabled"),
}

func runAutoUpdateEnableOrDisable(enabled bool, confirmMsg, successMsg string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client := cliutils.PackageClient(ctx)
		var pkgs []ctrlpkg.Package
		if autoUpdateEnabledDisabledOptions.All {
			if len(args) > 0 {
				fmt.Fprintf(os.Stderr, "Too many arguments: %v\n", args)
				cliutils.ExitWithError()
			}
			if autoUpdateEnabledDisabledOptions.Kind != KindPackage && autoUpdateEnabledDisabledOptions.Namespace == "" {
				var pkgList v1alpha1.ClusterPackageList
				if err := client.ClusterPackages().GetAll(ctx, &pkgList); err != nil {
					fmt.Fprintf(os.Stderr, "Could not list packages: %v", err)
					cliutils.ExitWithError()
				}
				for i := range pkgList.Items {
					pkgs = append(pkgs, &pkgList.Items[i])
				}
			}
			if autoUpdateEnabledDisabledOptions.Kind != KindClusterPackage {
				var pkgList v1alpha1.PackageList
				if err := client.Packages(autoUpdateEnabledDisabledOptions.Namespace).
					GetAll(ctx, &pkgList); err != nil {
					fmt.Fprintf(os.Stderr, "Could not list packages: %v", err)
					cliutils.ExitWithError()
				}
				for i := range pkgList.Items {
					pkgs = append(pkgs, &pkgList.Items[i])
				}
			}
			for _, pkg := range pkgs {
				args = append(args, pkg.GetName())
			}
		} else {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Please specify at least one package")
				cliutils.ExitWithError()
			}
			pkgs = make([]ctrlpkg.Package, len(args))
			for i, name := range args {
				pkg, err := getPackageOrClusterPackage(ctx, name,
					autoUpdateEnabledDisabledOptions.KindOptions,
					autoUpdateEnabledDisabledOptions.NamespaceOptions)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not get package %v: %v", name, err)
					cliutils.ExitWithError()
				}
				pkgs[i] = pkg
			}
		}

		if len(pkgs) == 0 {
			fmt.Fprintln(os.Stderr, "Nothing to do")
			cliutils.ExitSuccess()
		}

		fmt.Fprintln(os.Stderr, confirmMsg)
		for _, pkg := range pkgs {
			if pkg.IsNamespaceScoped() {
				fmt.Fprintf(os.Stderr, " * %v (Package in namespace %v with type %v)\n",
					pkg.GetName(), pkg.GetNamespace(), pkg.GetSpec().PackageInfo.Name)
			} else {
				fmt.Fprintf(os.Stderr, " * %v (ClusterPackage)\n", pkg.GetName())
			}
		}
		if !autoUpdateEnabledDisabledOptions.Yes && !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}

		var err error
		for _, pkg := range pkgs {
			if pkg.AutoUpdatesEnabled() != enabled {
				pkg.SetAutoUpdatesEnabled(enabled)
				switch pkg := pkg.(type) {
				case *v1alpha1.ClusterPackage:
					multierr.AppendInto(&err, client.ClusterPackages().Update(ctx, pkg))
				case *v1alpha1.Package:
					multierr.AppendInto(&err, client.Packages(pkg.Namespace).Update(ctx, pkg))
				default:
					panic("unexpected type")
				}
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error modifying some packages: %v", err)
			cliutils.ExitWithError()
		}
		fmt.Fprintf(os.Stderr, "%v: %v\n", successMsg, strings.Join(args, ", "))
		cliutils.ExitSuccess()
	}
}

var autoUpdateCmd = &cobra.Command{
	Use:   "auto-update",
	Short: "Update autopilot for packages where automatic updates are enabled",
	Args:  cobra.NoArgs,
	PreRun: cliutils.RunAll(
		func(c *cobra.Command, s []string) { config.NonInteractive = true },
		cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	),
	Run: runAutoUpdate,
}

func runAutoUpdate(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	client := cliutils.PackageClient(ctx)
	updater := update.NewUpdater(ctx).
		WithStatusWriter(statuswriter.Stderr())

	var pkgs []ctrlpkg.Package

	var cpkgList v1alpha1.ClusterPackageList
	if err := client.ClusterPackages().GetAll(ctx, &cpkgList); err != nil {
		panic(err)
	}

	for i, pkg := range cpkgList.Items {
		if pkg.AutoUpdatesEnabled() {
			pkgs = append(pkgs, &cpkgList.Items[i])
		}
	}

	var pkgList v1alpha1.PackageList
	if err := client.Packages("").GetAll(ctx, &pkgList); err != nil {
		panic(err)
	}

	for i, pkg := range pkgList.Items {
		if pkg.AutoUpdatesEnabled() {
			pkgs = append(pkgs, &pkgList.Items[i])
		}
	}

	if len(pkgs) == 0 {
		fmt.Fprintln(os.Stderr, "Automatic updates must be enabled for at least one package")
		cliutils.ExitSuccess()
	}

	tx, err := updater.Prepare(ctx, update.GetExact(pkgs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error preparing update: %v\n", err)
		cliutils.ExitWithError()
	}
	printTransaction(*tx)

	if updated, err := updater.ApplyBlocking(ctx, tx); err != nil {
		fmt.Fprintf(os.Stderr, "Error applying update: %v\n", err)
		cliutils.ExitWithError()
	} else {
		updatedNames := make([]string, len(updated))
		for i := range updated {
			updatedNames[i] = updated[i].GetName()
		}
		fmt.Fprintf(os.Stderr, "Updated packages: %v\n", strings.Join(updatedNames, ", "))
	}

	cliutils.ExitSuccess()
}

func init() {
	for _, cmd := range []*cobra.Command{autoUpdateEnableCmd, autoUpdateDisableCmd} {
		cmd.Flags().BoolVar(&autoUpdateEnabledDisabledOptions.Yes, "yes",
			autoUpdateEnabledDisabledOptions.Yes, "do not ask for confirmation")
		cmd.Flags().BoolVar(&autoUpdateEnabledDisabledOptions.All, "all",
			autoUpdateEnabledDisabledOptions.All, "set for all packages")
		autoUpdateEnabledDisabledOptions.KindOptions.AddFlagsToCommand(cmd)
		autoUpdateEnabledDisabledOptions.NamespaceOptions.AddFlagsToCommand(cmd)
		autoUpdateCmd.AddCommand(cmd)
	}
	RootCmd.AddCommand(autoUpdateCmd)
}
