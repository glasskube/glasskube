package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

var PackageGVR = v1alpha1.GroupVersion.WithResource("packages")

func InitKubeClient(kubeconfig string) (*PackageV1Alpha1Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		loadingRules.ExplicitPath = kubeconfig
	}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	err = v1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}
	pkgClient, err := NewPackageClient(config)
	if err != nil {
		return nil, err
	}
	return pkgClient, nil
}
