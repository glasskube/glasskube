package client

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	clientContextKey = iota
)

var PackageGVR = v1alpha1.GroupVersion.WithResource("packages")

func InitKubeConfig(kubeconfig string) (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		loadingRules.ExplicitPath = kubeconfig
	}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	return clientConfig.ClientConfig()
}

func SetupContext(ctx context.Context, kubeconfig string) (context.Context, error) {
	config, err := InitKubeConfig(kubeconfig)
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
	return context.WithValue(ctx, clientContextKey, pkgClient), nil
}

func FromContext(ctx context.Context) *PackageV1Alpha1Client {
	value := ctx.Value(clientContextKey)
	if value != nil {
		if client, ok := value.(*PackageV1Alpha1Client); ok {
			return client
		}
	}
	return nil
}
