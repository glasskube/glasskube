package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	pkgClient "github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/glasskube/glasskube/pkg/uninstall"
	"github.com/spf13/cobra"
)

var uninstallCmdOptions = struct {
	ForceUninstall bool
}{}

var uninstallCmd = &cobra.Command{
	Use:    "uninstall [package-name]",
	Short:  "Uninstall a package",
	Long:   `Uninstall a package.`,
	Args:   cobra.ExactArgs(1),
	PreRun: cliutils.SetupClientContext(true),
	Run: func(cmd *cobra.Command, args []string) {
		client := pkgClient.FromContext(cmd.Context())
		pkgName := args[0]
		var pkg v1alpha1.Package
		if err := client.Packages().Get(cmd.Context(), pkgName, &pkg); err != nil {
			fmt.Fprintf(os.Stderr, "Could not get installed package %v:\n%v\n", pkgName, err)
			os.Exit(1)
			return
		}
		proceed := uninstallCmdOptions.ForceUninstall || cliutils.YesNoPrompt(
			fmt.Sprintf(
				"%v will be removed from your cluster (%v). Are you sure?",
				pkgName,
				pkgClient.RawConfigFromContext(cmd.Context()).CurrentContext,
			),
			false,
		)
		if proceed {
			fmt.Printf("Uninstalling %v...\n", pkgName)
			if err := uninstall.NewUninstaller(client, &pkg).
				WithStatusWriter(statuswriter.Spinner()).
				UninstallBlocking(cmd.Context(), pkgName); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred during uninstallation:\n\n%v\n", err)
				os.Exit(1)
				return
			}
			fmt.Println("🗑️ Uninstalled successfully.")
		} else {
			fmt.Println("❌ Uninstallation cancelled.")
		}
	},
}

func init() {
	uninstallCmd.PersistentFlags().BoolVar(&uninstallCmdOptions.ForceUninstall, "force", false,
		"skip the confirmation question and uninstall right away")
	RootCmd.AddCommand(uninstallCmd)
}
