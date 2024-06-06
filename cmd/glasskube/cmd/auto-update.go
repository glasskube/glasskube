package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/update"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

var autoUpdateEnabledDisabledOptions = struct{ Yes, All bool }{}

var autoUpdateEnableCmd = &cobra.Command{
	Use:               "enable [...package]",
	Short:             "enable automatic updates for packages",
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeInstalledPackageNames,
	Run: runAutoUpdateEnableOrDisable(true,
		"Enable automatic updates for the following packages", "Automatic updates enabled"),
}

var autoUpdateDisableCmd = &cobra.Command{
	Use:               "disable [...package]",
	Short:             "disable automatic updates for packages",
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeInstalledPackageNames,
	Run: runAutoUpdateEnableOrDisable(false,
		"Enable automatic updates for the following packages", "Automatic updates disabled"),
}

func runAutoUpdateEnableOrDisable(enabled bool, confirmMsg, successMsg string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client := cliutils.PackageClient(ctx)
		var pkgs []v1alpha1.Package
		if autoUpdateEnabledDisabledOptions.All {
			if len(args) > 0 {
				fmt.Fprintf(os.Stderr, "Too many arguments: %v\n", args)
				cliutils.ExitWithError()
			}
			var pkgList v1alpha1.PackageList
			if err := client.Packages().GetAll(ctx, &pkgList); err != nil {
				fmt.Fprintf(os.Stderr, "Could not list packages: %v", err)
				cliutils.ExitWithError()
			}
			pkgs = pkgList.Items
			for _, pkg := range pkgs {
				args = append(args, pkg.Name)
			}
		} else {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Please specify at least one package")
				cliutils.ExitWithError()
			}
			pkgs = make([]v1alpha1.Package, len(args))
			for i, name := range args {
				var pkg v1alpha1.Package
				if err := client.Packages().Get(ctx, name, &pkg); err != nil {
					fmt.Fprintf(os.Stderr, "Could not get package %v: %v", name, err)
					cliutils.ExitWithError()
				}
				pkgs[i] = pkg
			}
		}

		fmt.Fprintf(os.Stderr, "%v: %v\n", confirmMsg, strings.Join(args, ", "))
		if !autoUpdateEnabledDisabledOptions.Yes && !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}

		var err error
		for _, pkg := range pkgs {
			if pkg.AutoUpdatesEnabled() != enabled {
				pkg.SetAutoUpdatesEnabled(enabled)
				multierr.AppendInto(&err, client.Packages().Update(ctx, &pkg))
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
	Short: "update autopilote for packages where automatic updates are enabled",
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

	var pkgs v1alpha1.PackageList
	if err := client.Packages().GetAll(ctx, &pkgs); err != nil {
		panic(err)
	}

	var packageNames []string
	for _, pkg := range pkgs.Items {
		if pkg.AutoUpdatesEnabled() {
			packageNames = append(packageNames, pkg.Name)
		}
	}
	if len(packageNames) == 0 {
		fmt.Fprintln(os.Stderr, "Automatic updates must be enabled for at least one package")
		cliutils.ExitSuccess()
	}

	tx, err := updater.Prepare(ctx, packageNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error preparing update: %v\n", err)
		cliutils.ExitWithError()
	}
	printTransaction(*tx)

	if updated, err := updater.Apply(ctx, tx); err != nil {
		fmt.Fprintf(os.Stderr, "Error applying update: %v\n", err)
		cliutils.ExitWithError()
	} else {
		updatedNames := make([]string, len(updated))
		for i := range updated {
			updatedNames[i] = updated[i].Name
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
		autoUpdateCmd.AddCommand(cmd)
	}
	RootCmd.AddCommand(autoUpdateCmd)
}
