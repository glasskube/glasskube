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
	return postProcess(clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil))
}

func FromBytes(data []byte) (*rest.Config, *api.Config, error) {
	if config, err := clientcmd.NewClientConfigFromBytes(data); err != nil {
		return nil, nil, err
	} else {
		return postProcess(config)
	}
}

func postProcess(clientConfig clientcmd.ClientConfig) (*rest.Config, *api.Config, error) {
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, nil, err
	}
	restConfig.APIPath = "/api"
	restConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	_ = rest.SetKubernetesDefaults(restConfig)
	return restConfig, &rawConfig, nil
}
