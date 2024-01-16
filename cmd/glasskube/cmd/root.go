package cmd

import (
	"github.com/glasskube/glasskube/cmd/glasskube/config"

	"github.com/spf13/cobra"
)

var (
	RootCmd = cobra.Command{
		Use:     "glasskube",
		Version: "0.0.0",
		Short:   "Kubernetes Package Management the easy way 🔥",
	}
)

func init() {
	RootCmd.PersistentFlags().StringVar(&config.Kubeconfig, "kubeconfig", "",
		"path to the kubeconfig file, whose current-context will be used (defaults to ~/.kube/config)")
}
