package cmd

import (
	"fmt"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/spf13/cobra"
)

type outputFormat string

const (
	outputFormatJSON outputFormat = outputFormat(clientutils.OutputFormatJSON)
	outputFormatYAML outputFormat = outputFormat(clientutils.OutputFormatYAML)
)

func (of *outputFormat) String() string {
	return string(*of)
}

func (of *outputFormat) Set(value string) error {
	switch value {
	case string(outputFormatJSON), string(outputFormatYAML):
		*of = outputFormat(value)
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", value)
	}
}

func (of *outputFormat) Type() string {
	return fmt.Sprintf("(%v|%v)", outputFormatJSON, outputFormatYAML)
}

func (of *outputFormat) OutputFormat() clientutils.OutputFormat {
	return clientutils.OutputFormat(of.String())
}

type OutputOptions struct {
	Output  outputFormat
	ShowAll bool
}

func (opts *OutputOptions) AddFlagsToCommand(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.VarP(&opts.Output, "output", "o", "Output format")
	flags.BoolVar(&opts.ShowAll, "show-all", false, "Show the complete output if -o is given")
}
