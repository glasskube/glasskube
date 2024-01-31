package cliutils

import (
	"fmt"
	"github.com/glasskube/glasskube/internal/bootstrap"
	"os"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/spf13/cobra"
)

func SetupClientContext(cmd *cobra.Command, args []string) {
	cfg := RequireConfig(config.Kubeconfig)
	bootstrap.RequireBootstrapped(cmd.Context(), cfg)
	if ctx, err := client.SetupContext(cmd.Context(), cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up the client:\n\n%v\n", err)
		os.Exit(1)
	} else {
		cmd.SetContext(ctx)
	}
}
