package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/web"
	"github.com/spf13/cobra"
)

var serveCmdOptions struct {
	port int
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start", "ui"},
	Short:   "Open UI",
	Long:    `Start server and open the UI.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		options := web.ServerOptions{
			Host:       "localhost",
			Port:       int32(serveCmdOptions.port),
			Kubeconfig: config.Kubeconfig,
		}
		server := web.NewServer(options)
		if err := server.Start(cmd.Context()); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred starting the webserver:\n\n%v\n", err)
			cliutils.ExitWithError(cmd.Context())
		}
	},
}

func init() {
	serveCmd.Flags().IntVarP(&serveCmdOptions.port, "port", "p", 8580, "Port for the webserver")
	RootCmd.AddCommand(serveCmd)
}
