package cmd

import (
	"github.com/spf13/cobra"
)

type DryRunOptions struct {
	DryRun bool
}

func (opt *DryRunOptions) AddFlagsToCommand(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&opt.DryRun, "dry-run", false, "Simulate the execution of the command without making any changes")
}
