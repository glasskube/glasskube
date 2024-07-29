package cmd

import (
	"fmt"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/spf13/cobra"
)

type outputFormatValue string

const (
	// TODO make better
	outputFormatJSON outputFormatValue = outputFormatValue(clientutils.OutputFormatJSON)
	outputFormatYAML outputFormatValue = outputFormatValue(clientutils.OutputFormatYAML)
)

func (of *outputFormatValue) String() string {
	return string(*of)
}

func (of *outputFormatValue) Set(value string) error {
	switch value {
	case string(outputFormatJSON), string(outputFormatYAML):
		*of = outputFormatValue(value)
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", value)
	}
}

func (of *outputFormatValue) Type() string {
	return fmt.Sprintf("(%v|%v)", outputFormatJSON, outputFormatYAML)
}

func (of *outputFormatValue) OutputFormat() clientutils.OutputFormat {
	return clientutils.OutputFormat(of.String())
}

type OutputOptions struct {
	Output  outputFormatValue
	ShowAll bool
}

func (opts *OutputOptions) AddFlagsToCommand(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.VarP(&opts.Output, "output", "o", "Output format")
	flags.BoolVar(&opts.ShowAll, "show-all", false,
		"Additionally print metadata fields other than name, namespace, annotations and labels")
}
