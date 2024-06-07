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
	port        int
	logLevel    int
	skipOpenFlg bool
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start", "ui"},
	Short:   "Open UI",
	Long:    `Start server and open the UI.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		options := web.ServerOptions{
			Host:               "localhost",
			Port:               int32(serveCmdOptions.port),
			Kubeconfig:         config.Kubeconfig,
			LogLevel:           serveCmdOptions.logLevel,
			SkipOpeningBrowser: serveCmdOptions.skipOpenFlg,
		}
		server := web.NewServer(options)
		if err := server.Start(cmd.Context()); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred starting the webserver:\n\n%v\n", err)
			cliutils.ExitWithError()
		}
	},
}

func init() {
	serveCmd.Flags().IntVarP(&serveCmdOptions.port, "port", "p", 8580, "Port for the webserver")
	serveCmd.Flags().IntVarP(&serveCmdOptions.logLevel, "log-level", "l", 0,
		"Level for additional logging, where 0 is the least verbose")
	serveCmd.Flags().BoolVarP(&serveCmdOptions.skipOpenFlg, "skip-open", "s", false, "Skip opening the browser")
	RootCmd.AddCommand(serveCmd)
}
