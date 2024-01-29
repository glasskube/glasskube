package cliutils

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/cmd/glasskube/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/spf13/cobra"
)

func SetupClientContext(cmd *cobra.Command, args []string) {
	if ctx, err := client.SetupContext(cmd.Context(), RequireConfig(config.Kubeconfig)); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up the client:\n\n%v\n", err)
		os.Exit(1)
	} else {
		cmd.SetContext(ctx)
	}
}
