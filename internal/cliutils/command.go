package cliutils

import "github.com/spf13/cobra"

func RunAll(funcs ...func(*cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		for _, f := range funcs {
			f(cmd, args)
		}
	}
}
