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

type OpenCmdOptions struct {
	NamespaceOptions
	KindOptions
	Port int32
}

var (
	openCmdOptions = OpenCmdOptions{
		KindOptions: DefaultKindOptions(),
	}
)

var openCmd = &cobra.Command{
	Use:   "open <package-name> [<entrypoint>]",
	Short: "Open the Web UI of a package",
	Long: `Open the Web UI of a package.
If the package manifest has more than one entrypoint, specify the name of the entrypoint to open.`,
	Args:   cobra.RangeArgs(1, 2),
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		pkgName := args[0]
		entrypointName := ""
		if len(args) == 2 {
			entrypointName = args[1]
		}

		pkg, err := getPackageOrClusterPackage(ctx, pkgName, openCmdOptions.KindOptions, openCmdOptions.NamespaceOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not get resource %v: %v\n", pkgName, err)
			cliutils.ExitWithError()
		}

		result, err := open.NewOpener().Open(ctx, pkg, entrypointName, openCmdOptions.Port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not open package %v: %v\n", pkgName, err)
			cliutils.ExitWithError()
		}

		stopBeforeExit := func() {
			fmt.Fprintln(os.Stderr, "🛑 Terminating forwarders...")
			result.Stop()
			cliutils.ExitFromSignal(nil)
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
					cliutils.ExitWithError()
				} else {
					fmt.Fprintln(os.Stderr, "❗ Forwarders closed unexpectedly")
					cliutils.ExitWithError()
				}
			}
		}
	},
}

func init() {
	openCmdOptions.KindOptions.AddFlagsToCommand(openCmd)
	openCmdOptions.NamespaceOptions.AddFlagsToCommand(openCmd)
	openCmd.Flags().Int32Var(&openCmdOptions.Port, "port", openCmdOptions.Port, "custom port for opening the package")
	RootCmd.AddCommand(openCmd)
}
