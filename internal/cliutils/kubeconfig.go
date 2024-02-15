package cliutils

import (
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"k8s.io/client-go/rest"
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

func RequireConfig(filePath string) (*rest.Config, *api.Config) {
	if config, rawConfig, err := kubeconfig.New(filePath); err != nil {
		if clientcmd.IsEmptyConfig(err) {
			fmt.Fprintln(os.Stderr, helpEmptyConfig)
		} else {
			fmt.Fprintf(os.Stderr, "Your kubeconfig file is invalid:\n\n%v\n", err)
		}
		os.Exit(1)
		return nil, nil
	} else {
		return config, rawConfig
	}
}
