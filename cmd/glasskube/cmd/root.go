package cmd

import (
	"fmt"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

var rootCmdOptions struct {
	SkipUpdateCheck bool
}

var (
	RootCmd = cobra.Command{
		Use:     "glasskube",
		Version: config.Version,
		Short:   "🧊 The missing Package Manager for Kubernetes 📦",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !rootCmdOptions.SkipUpdateCheck {
				cliutils.UpdateFetch()
			}
		},
	}
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&rootCmdOptions.SkipUpdateCheck, "skip-update-check", config.IsDevBuild(),
		"Do not check for Glasskube updates")
	RootCmd.PersistentFlags().StringVar(&config.Kubeconfig, "kubeconfig", "",
		fmt.Sprintf("path to the kubeconfig file, whose current-context will be used (defaults to %v)",
			clientcmd.RecommendedHomeFile))
}
