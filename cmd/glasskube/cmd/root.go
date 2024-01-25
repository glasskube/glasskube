package cmd

import (
	"fmt"

	"github.com/glasskube/glasskube/cmd/glasskube/config"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

var (
	RootCmd = cobra.Command{
		Use:     "glasskube",
		Version: "0.0.0",
		Short:   "🧊 The missing Package Manager for Kubernetes 📦",
	}
)

func init() {
	RootCmd.PersistentFlags().StringVar(&config.Kubeconfig, "kubeconfig", "",
		fmt.Sprintf("path to the kubeconfig file, whose current-context will be used (defaults to %v)", clientcmd.RecommendedHomeFile))
}
