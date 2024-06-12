package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/glasskube/glasskube/pkg/purge"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/spf13/cobra"
)

type purgeOptions struct {
	yes bool
}

var purgeCmdOptions = purgeOptions{}

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge Glasskube from a Kubernetes cluster",
	Long: "Purges Glasskube from a Kubernetes cluster, " +
		"thereby uninstalling the Glasskube operator and all related resources.",
	Args:   cobra.NoArgs,
	PreRun: cliutils.SetupClientContext(false, util.Pointer(true)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)
		client := purge.NewPurger(cfg)
		ctx := cmd.Context()

		client.WithStatusWriter(statuswriter.Spinner())

		currentContext := clicontext.RawConfigFromContext(ctx).CurrentContext

		if !purgeCmdOptions.yes {
			confirmMessage := fmt.Sprintf("Glasskube will be purged from context %s.\nContinue? ", currentContext)
			if !cliutils.YesNoPrompt(confirmMessage, true) {
				fmt.Println("Operation stopped")
				cliutils.ExitWithError()
			}
		}

		IsBootstrapped, err := bootstrap.IsBootstrapped(cmd.Context(), cfg)
		if err != nil && !IsBootstrapped {
			fmt.Printf("error : %v\n", err)
			cliutils.ExitWithError()
		}

		if err := client.Purge(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during purge:\n%v\n", err)
			cliutils.ExitWithError()
		}

	},
}

func init() {
	RootCmd.AddCommand(purgeCmd)
	purgeCmd.Flags().BoolVar(&purgeCmdOptions.yes, "yes", false, "Skip confirmation prompt")
}
