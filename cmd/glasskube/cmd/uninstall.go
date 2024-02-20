package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	pkgClient "github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
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
		pkg, err := list.Get(client, cmd.Context(), pkgName)
		if err != nil {
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
			fmt.Printf("Uninstalling %v.\n", pkgName)
			err = uninstall.Uninstall(client, cmd.Context(), pkg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred during uninstallation:\n\n%v\n", err)
				os.Exit(1)
				return
			}
			fmt.Println("Uninstalled successfully.")
		} else {
			fmt.Println("Uninstallation cancelled.")
		}
	},
}

func init() {
	uninstallCmd.PersistentFlags().BoolVar(&uninstallCmdOptions.ForceUninstall, "force", false,
		"skip the confirmation question and uninstall right away")
	RootCmd.AddCommand(uninstallCmd)
}
