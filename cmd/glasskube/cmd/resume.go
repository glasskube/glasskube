package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/suspend"
	"github.com/spf13/cobra"
)

var resumeCmdOptions = struct {
	KindOptions
	NamespaceOptions
}{
	KindOptions: DefaultKindOptions(),
}

var resumeCmd = &cobra.Command{
	Use:               "resume <package-name>",
	Short:             "Resume reconciliation of a previously suspended package",
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run:               func(cmd *cobra.Command, args []string) { runResume(cmd.Context(), args[0]) },
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: installedPackagesCompletionFunc(&resumeCmdOptions.NamespaceOptions, &resumeCmdOptions.KindOptions),
}

func runResume(ctx context.Context, name string) {
	pkg, err := getPackageOrClusterPackage(ctx, name, resumeCmdOptions.KindOptions,
		resumeCmdOptions.NamespaceOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		cliutils.ExitWithError()
	}

	if resumed, err := suspend.Resume(ctx, pkg); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		cliutils.ExitWithError()
	} else if resumed {
		fmt.Fprintf(os.Stderr, "✅ %v is now resumed\n", pkg.GetName())
	} else {
		fmt.Fprintf(os.Stderr, "☑️  %v was not suspended\n", pkg.GetName())
	}

	cliutils.ExitSuccess()
}

func init() {
	resumeCmdOptions.KindOptions.AddFlagsToCommand(resumeCmd)
	resumeCmdOptions.NamespaceOptions.AddFlagsToCommand(resumeCmd)
	RootCmd.AddCommand(resumeCmd)
}
