package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/telemetry"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

var rootCmdOptions struct {
	SkipUpdateCheck bool
	NoProgress      bool
}

var (
	RootCmd = cobra.Command{
		Use:     "glasskube",
		Version: config.Version,
		Short:   "🧊 The next generation Package Manager for Kubernetes 📦",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			telemetry.Init()
			if !rootCmdOptions.SkipUpdateCheck {
				cliutils.UpdateFetch()
			}

			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)
			go func() {
				sig := <-signals
				// TODO find another way
				if cmd.Name() != openCmd.Name() {
					cliutils.ExitFromSignal(&sig)
				}
			}()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cliutils.ExitSuccess()
		},
	}
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&rootCmdOptions.SkipUpdateCheck, "skip-update-check", config.IsDevBuild(),
		"Do not check for Glasskube updates")
	RootCmd.PersistentFlags().StringVar(&config.Kubeconfig, "kubeconfig", "",
		fmt.Sprintf("Path to the kubeconfig file, whose current-context will be used (defaults to %v)",
			clientcmd.RecommendedHomeFile))
	RootCmd.PersistentFlags().BoolVar(&config.NonInteractive, "non-interactive", config.NonInteractive,
		"Run in non-interactive mode. "+
			"If interactivity would be required, the command will terminate with a non-zero exit code.")
	RootCmd.PersistentFlags().BoolVar(&rootCmdOptions.NoProgress, "no-progress", false,
		"Prevent progress logging to the cli")
}
