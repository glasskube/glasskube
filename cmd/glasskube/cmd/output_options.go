package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type OutputFormat string

const (
	OutputFormatJSON OutputFormat = "json"
	OutputFormatYAML OutputFormat = "yaml"
)

func (of *OutputFormat) String() string {
	return string(*of)
}

func (of *OutputFormat) Set(value string) error {
	switch value {
	case string(OutputFormatJSON), string(OutputFormatYAML):
		*of = OutputFormat(value)
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", value)
	}
}

func (of *OutputFormat) Type() string {
	return fmt.Sprintf("(%v|%v)", OutputFormatJSON, OutputFormatYAML)
}

type OutputOptions struct {
	Output OutputFormat
}

func (opts *OutputOptions) AddFlagsToCommand(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.VarP(&opts.Output, "output", "o", "output format")
}
