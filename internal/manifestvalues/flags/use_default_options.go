package flags

import "github.com/spf13/cobra"

type UseDefaultOptions struct {
	UseDefault []string
}

func (opts *UseDefaultOptions) AddFlagsToCommand(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringArrayVar(&opts.UseDefault, "use-default", opts.UseDefault,
		"Instruct glasskube to use the default value for the speciefied definition name(s).\n"+
			"Specify \"all\" to use all available default values.")
}
