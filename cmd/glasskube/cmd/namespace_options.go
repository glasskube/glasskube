package cmd

import (
	"context"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/spf13/cobra"
)

type NamespaceOptions struct {
	Namespace string
}

func (opt *NamespaceOptions) AddFlagsToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opt.Namespace, "namespace", "n", opt.Namespace, "namespace for resources")
}

func (opt *NamespaceOptions) GetActualNamespace(ctx context.Context) string {
	if opt.Namespace != "" {
		return opt.Namespace
	} else {
		rawConfig := clicontext.RawConfigFromContext(ctx)
		if current, ok := rawConfig.Contexts[rawConfig.CurrentContext]; ok && current.Namespace != "" {
			return current.Namespace
		} else {
			return "default"
		}
	}
}
