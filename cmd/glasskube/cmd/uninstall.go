package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/uninstall"
	"github.com/spf13/cobra"
)

var uninstallCmdOptions = struct {
	NoWait          bool
	Yes             bool
	DeleteNamespace bool
	KindOptions
	NamespaceOptions
}{
	KindOptions: DefaultKindOptions(),
}

var uninstallCmd = &cobra.Command{
	Use:               "uninstall <package-name>",
	Short:             "Uninstall a package",
	Long:              `Uninstall a package.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeInstalledPackageNames,
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		ctx := cmd.Context()
		currentContext := clicontext.RawConfigFromContext(ctx).CurrentContext
		client := clicontext.PackageClientFromContext(ctx)
		dm := cliutils.DependencyManager(ctx)
		uninstaller := uninstall.NewUninstaller(client)
		if !rootCmdOptions.NoProgress {
			uninstaller.WithStatusWriter(statuswriter.Spinner())
		}

		pkg, err := getPackageOrClusterPackage(
			ctx, pkgName, uninstallCmdOptions.KindOptions, uninstallCmdOptions.NamespaceOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not get resource: %v\n", err)
			cliutils.ExitWithError()
		}

		var deleteNamespace bool
		if !pkg.IsNamespaceScoped() {
			if g, err := dm.NewGraph(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error validating uninstall: %v\n", err)
				cliutils.ExitWithError()
			} else {
				g.Delete(pkgName)
				pruned := g.Prune()
				if err := g.Validate(); err != nil {
					fmt.Fprintf(os.Stderr, "❌ %v can not be uninstalled for the following reason: %v\n", pkgName, err)
					cliutils.ExitWithError()
				} else {
					showUninstallDetails(currentContext, pkgName, pruned)
					if !uninstallCmdOptions.Yes && !cliutils.YesNoPrompt("Do you want to continue?", false) {
						fmt.Println("❌ Uninstallation cancelled.")
						cliutils.ExitSuccess()
					}
				}
			}
		} else if uninstallCmdOptions.DeleteNamespace {
			// list all packages in the namespace
			lister := list.NewListerWithRepoCache(ctx)
			namespace := pkg.GetNamespace()
			pkgsInNamespace, err := lister.GetPackagesWithStatus(ctx, list.ListOptions{
				Namespace: namespace,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Could not list packages: %v\n", err)
				cliutils.ExitWithError()
			}
			pkgCount := 0
			for _, p := range pkgsInNamespace {
				pkgCount += len(p.Packages)
			}
			if pkgCount > 1 {
				fmt.Fprintf(os.Stderr, "❌ Namespace %v contains more than one package. "+
					"Current count of packages present in the namespace: %v\n", namespace, pkgCount)
				cliutils.ExitWithError()
			}
			deleteNamespace = true
		}

		// delete the namespace only if it is namespace scoped and there are no other packages in the namespace since
		// cluster packages delete the namespace once they are deleted
		if uninstallCmdOptions.NoWait {
			if err := uninstaller.UninstallAndDeleteNamespace(ctx, pkg, deleteNamespace); err != nil {
				fmt.Fprintf(os.Stderr, "\n❌ An error occurred during uninstallation:\n\n%v\n", err)
				cliutils.ExitWithError()
			}
			fmt.Fprintln(os.Stderr, "Uninstallation started in background")
		} else {
			if err := uninstaller.UninstallAndDeleteNamespaceBlocking(ctx, pkg, deleteNamespace); err != nil {
				fmt.Fprintf(os.Stderr, "\n❌ An error occurred during uninstallation:\n\n%v\n", err)
				cliutils.ExitWithError()
			}
			fmt.Fprintf(os.Stderr, "🗑️  %v uninstalled successfully.\n", pkgName)
		}
	},
}

func showUninstallDetails(context, name string, pruned []string) {
	fmt.Fprintf(os.Stderr,
		"The following packages will be %v from your cluster (%v):\n",
		color.New(color.Bold).Sprint("removed"),
		context)
	fmt.Fprintf(os.Stderr, " * %v (requested by user)\n", name)
	for _, dep := range pruned {
		fmt.Fprintf(os.Stderr, " * %v (dependency)\n", dep)
	}
}

func init() {
	uninstallCmdOptions.KindOptions.AddFlagsToCommand(uninstallCmd)
	uninstallCmdOptions.NamespaceOptions.AddFlagsToCommand(uninstallCmd)
	uninstallCmd.PersistentFlags().BoolVarP(&uninstallCmdOptions.DeleteNamespace, "delete-namespace", "d", false,
		"Delete the namespace if it is empty after deleting the resources")
	uninstallCmd.PersistentFlags().BoolVar(&uninstallCmdOptions.NoWait, "no-wait", false,
		"Perform non-blocking uninstall")
	uninstallCmd.PersistentFlags().BoolVarP(&uninstallCmdOptions.Yes, "yes", "y", false,
		"Do not ask for any confirmation")
	RootCmd.AddCommand(uninstallCmd)
}
