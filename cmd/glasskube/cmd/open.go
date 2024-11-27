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
	host string
	port int32
}

var (
	openCmdOptions = OpenCmdOptions{
		KindOptions: DefaultKindOptions(),
		host:        "localhost",
	}
)

var openCmd = &cobra.Command{
	Use:   "open <package-name> [<entrypoint>]",
	Short: "Open the Web UI of a package",
	Long: `Open the Web UI of a package.
If the package manifest has more than one entrypoint, specify the name of the entrypoint to open.`,
	Args: cobra.RangeArgs(1, 2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			delegate := installedPackagesCompletionFunc(&openCmdOptions.NamespaceOptions, &openCmdOptions.KindOptions)
			return delegate(cmd, args, toComplete)
		} else {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	},
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

		result, err := open.NewOpener().Open(ctx, pkg, entrypointName, openCmdOptions.host, openCmdOptions.port)
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
	openCmd.Flags().StringVar(&openCmdOptions.host, "host", openCmdOptions.host,
		"Custom hostname to open the local port on")
	openCmd.Flags().Int32Var(&openCmdOptions.port, "port", openCmdOptions.port, "Custom port for opening the package")
	RootCmd.AddCommand(openCmd)
}
