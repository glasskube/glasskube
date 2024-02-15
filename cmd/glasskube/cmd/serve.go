package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/web"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var serveCmdOptions struct {
	port int
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start", "ui"},
	Short:   "Open UI",
	Long:    `Start server and open the UI.`,
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var support *web.ServerConfigSupport
		cfg, rawCfg, err := kubeconfig.New(config.Kubeconfig)
		if err != nil {
			support = &web.ServerConfigSupport{
				KubeconfigError:           err,
				KubeconfigDefaultLocation: clientcmd.RecommendedHomeFile,
			}
			if clientcmd.IsEmptyConfig(err) {
				support.KubeconfigMissing = true
			}
		}

		var ctx = cmd.Context()
		if cfg != nil {
			ctx, err = client.SetupContext(ctx, cfg, rawCfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred starting the webserver:\n\n%v\n", err)
				os.Exit(1)
				return
			} else {
				cmd.SetContext(ctx)
			}
		}

		server := web.NewServer("localhost", int32(serveCmdOptions.port))
		if err = server.Start(ctx, support); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred starting the webserver:\n\n%v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	serveCmd.Flags().IntVarP(&serveCmdOptions.port, "port", "p", 8580, "Port for the webserver")
	RootCmd.AddCommand(serveCmd)
}
