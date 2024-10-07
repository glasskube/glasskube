package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/suspend"
	"github.com/spf13/cobra"
)

var suspendCmdOptions = struct {
	KindOptions
	NamespaceOptions
}{
	KindOptions: DefaultKindOptions(),
}

var suspendCmd = &cobra.Command{
	Use:     "suspend <package-name>",
	Short:   "Suspend reconciliation of a package",
	Aliases: []string{"pause"},
	PreRun:  cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run:     func(cmd *cobra.Command, args []string) { runSuspend(cmd.Context(), args[0]) },
	Args:    cobra.ExactArgs(1),
	ValidArgsFunction: installedPackagesCompletionFunc(
		&suspendCmdOptions.NamespaceOptions,
		&suspendCmdOptions.KindOptions,
	),
}

func runSuspend(ctx context.Context, name string) {
	pkg, err := getPackageOrClusterPackage(ctx, name, suspendCmdOptions.KindOptions,
		suspendCmdOptions.NamespaceOptions)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cliutils.ExitWithError()
	}

	result, err := suspend.Suspend(ctx, pkg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cliutils.ExitWithError()
	}

	switch result {
	case suspend.Suspended:
		fmt.Fprintf(os.Stderr, "%v has been suspended\n", pkg.GetName())
	case suspend.UpToDate:
		fmt.Fprintf(os.Stderr, "%v is already suspended\n", pkg.GetName())
	default:
		fmt.Fprintf(os.Stderr, "unexpected suspend result: %v\n", result)
		cliutils.ExitWithError()
	}
	cliutils.ExitSuccess()
}

func init() {
	suspendCmdOptions.KindOptions.AddFlagsToCommand(suspendCmd)
	suspendCmdOptions.NamespaceOptions.AddFlagsToCommand(suspendCmd)
	RootCmd.AddCommand(suspendCmd)
}
