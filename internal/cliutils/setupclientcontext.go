package cliutils

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/cmd/glasskube/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	helpEmptyConfig = fmt.Sprintf(`
Your kubeconfig file is either empty or missing!
Please, download the kubeconfig file from your cloud provider and copy it to the default location, or use the --kubeconfig flag.
The default location is: %v

If you want to test glasskube locally, check out minikube: https://minikube.sigs.k8s.io/
	`, clientcmd.RecommendedHomeFile)
)

func SetupClientContext(cmd *cobra.Command, args []string) {
	ctx, err := client.SetupContext(cmd.Context(), config.Kubeconfig)
	if err != nil {
		if clientcmd.IsEmptyConfig(err) {
			fmt.Fprintln(os.Stderr, helpEmptyConfig)
		} else {
			fmt.Fprintf(os.Stderr, "Your kubeconfig file is invalid:\n\n%v\n", err)
		}
		os.Exit(1)
	} else {
		cmd.SetContext(ctx)
	}
}
