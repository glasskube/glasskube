package cmd

import (
	"github.com/spf13/cobra"
)

type OutputFormat string

const (
	OutputFormatJSON OutputFormat = "json"
	OutputFormatYAML OutputFormat = "yaml"
)

type OutputOptions struct {
	Output OutputFormat
}

func (opts *OutputOptions) AddFlagsToCommand(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVarP((*string)(&opts.Output), "output", "o", "", "Output format (json|yaml)")
}
