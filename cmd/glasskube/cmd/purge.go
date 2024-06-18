package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
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
	PreRun: cliutils.SetupClientContext(false, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)
		client := purge.NewPurger(cfg)
		ctx := cmd.Context()

		if !rootCmdOptions.NoProgress {
			client.WithStatusWriter(statuswriter.Spinner())
		}

		bold := color.New(color.Bold).SprintFunc()
		currentContext := clicontext.RawConfigFromContext(ctx).CurrentContext

		isBootstrapped, err := bootstrap.IsBootstrapped(cmd.Context(), cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			cliutils.ExitWithError()
		}
		if !isBootstrapped {
			fmt.Fprintln(os.Stderr, "error: glasskube is not bootstrapped")
			cliutils.ExitWithError()
		}

		if !purgeCmdOptions.yes {
			confirmMessage := fmt.Sprintf("⚠️ Glasskube and all related resources will be purged from context %s."+
				"\nThis includes removal of all installed packages!\nContinue? ", bold(currentContext))
			if !cliutils.YesNoPrompt(confirmMessage, true) {
				fmt.Fprintln(os.Stderr, "Operation stopped")
				cliutils.ExitWithError()
			}
		}

		if err := client.Purge(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during purge:\n%v\n", err)
			cliutils.ExitWithError()
		}

	},
}

func init() {
	RootCmd.AddCommand(purgeCmd)
	purgeCmd.Flags().BoolVar(&purgeCmdOptions.yes, "yes", false, "skip confirmation prompt")
}
