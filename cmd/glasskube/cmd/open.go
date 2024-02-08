package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/open"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [package-name] [entrypoint]",
	Short: "Open the Web UI of a package",
	Long: `Open the Web UI of a package.
If the package manifest has more than one entrypoint, specify the name of the entrypoint to open.`,
	Args:   cobra.RangeArgs(1, 2),
	PreRun: cliutils.SetupClientContext(true),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		var entrypointName string
		if len(args) == 2 {
			entrypointName = args[1]
		}

		result, err := open.NewOpener().Open(cmd.Context(), pkgName, entrypointName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not open package: %v\n", err)
			os.Exit(1)
		}

		stopBeforeExit := func() {
			fmt.Fprintln(os.Stderr, "🛑 Terminating forwarders...")
			result.Stop()
		}
		defer stopBeforeExit()

		stopCh := make(chan os.Signal, 1)
		signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

		go func() {
			result.WaitReady()
			fmt.Fprintf(os.Stderr, "✅ %s is now reachable at %s\n", pkgName, result.Url)
			if err = cliutils.OpenInBrowser(result.Url); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Could not open browser: %v\n", err)
			}
		}()

	outer:
		for {
			select {
			case <-stopCh:
				fmt.Fprintln(os.Stderr, "👋 Received interrupt signal")
				break outer
			case err := <-result.Completion:
				if err != nil {
					fmt.Fprintf(os.Stderr, "❌ An error occurred: %v\n", err)
					stopBeforeExit()
					os.Exit(1)
				} else {
					fmt.Fprintln(os.Stderr, "❗ Forwarders closed unexpectedly")
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
