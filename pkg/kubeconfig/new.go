package kubeconfig

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func New(filePath string) (*rest.Config, *api.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if filePath != "" {
		loadingRules.ExplicitPath = filePath
	}
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	config, err := loader.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	rawConfig, err := loader.RawConfig()
	if err != nil {
		return nil, nil, err
	}
	config.APIPath = "/api"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	_ = rest.SetKubernetesDefaults(config)
	return config, &rawConfig, nil
}
