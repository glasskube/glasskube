package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/glasskube/glasskube/internal/cliutils"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/web"
	"github.com/spf13/cobra"
)

type ServeCmdOptions struct {
	host     string
	port     int
	logLevel int
	skipOpen bool
}

func (opts ServeCmdOptions) ServerOptions() web.ServerOptions {
	return web.ServerOptions{
		Host:               opts.host,
		Port:               strconv.Itoa(opts.port),
		Kubeconfig:         config.Kubeconfig,
		LogLevel:           opts.logLevel,
		SkipOpeningBrowser: opts.skipOpen,
	}
}

var (
	serveCmdOptions = ServeCmdOptions{
		host: "localhost",
		port: 8050,
	}
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start", "ui"},
	Short:   "Open UI",
	Long:    `Start server and open the UI.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server := web.NewServer(serveCmdOptions.ServerOptions())
		if err := server.Start(cmd.Context()); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred starting the webserver:\n\n%v\n", err)
			cliutils.ExitWithError()
		}
	},
}

func init() {
	serveCmd.Flags().StringVar(&serveCmdOptions.host, "host", serveCmdOptions.host,
		"Hostname for the webserver")
	serveCmd.Flags().IntVarP(&serveCmdOptions.port, "port", "p", serveCmdOptions.port,
		"Port for the webserver")
	serveCmd.Flags().IntVarP(&serveCmdOptions.logLevel, "log-level", "l", serveCmdOptions.logLevel,
		"Level for additional logging, where 0 is the least verbose")
	serveCmd.Flags().BoolVarP(&serveCmdOptions.skipOpen, "skip-open", "s", serveCmdOptions.skipOpen,
		"Skip opening the browser")
	RootCmd.AddCommand(serveCmd)
}
