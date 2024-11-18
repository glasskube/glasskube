package clicontext

import (
	"context"

	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/pkg/client"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

type (
	ContextKey int
)

const (
	pkgClientContextKey ContextKey = iota
	k8sClientContextKey
	configContextKey
	rawConfigContextKey
	repoClientsetContextKey
)

func SetupContext(ctx context.Context, config *rest.Config, rawConfig *api.Config) (context.Context, error) {
	if pkgClient, err := client.New(config); err != nil {
		return nil, err
	} else if k8s, err := kubernetes.NewForConfig(config); err != nil {
		return nil, err
	} else {
		return SetupContextWithClient(ctx, config, rawConfig, pkgClient, k8s), nil
	}
}

func SetupContextWithClient(
	ctx context.Context,
	config *rest.Config,
	rawConfig *api.Config,
	client client.PackageV1Alpha1Client,
	k8s *kubernetes.Clientset,
) context.Context {
	ctx = context.WithValue(ctx, pkgClientContextKey, client)
	ctx = context.WithValue(ctx, k8sClientContextKey, k8s)
	ctx = context.WithValue(ctx, configContextKey, config)
	ctx = context.WithValue(ctx, rawConfigContextKey, rawConfig)
	return ctx
}

func ContextWithRepositoryClientset(parent context.Context, clientset repoclient.RepoClientset) context.Context {
	return context.WithValue(parent, repoClientsetContextKey, clientset)
}

func PackageClientFromContext(ctx context.Context) client.PackageV1Alpha1Client {
	value := ctx.Value(pkgClientContextKey)
	if value != nil {
		if client, ok := value.(client.PackageV1Alpha1Client); ok {
			return client
		}
	}
	return nil
}

func KubernetesClientFromContext(ctx context.Context) *kubernetes.Clientset {
	value := ctx.Value(k8sClientContextKey)
	if value != nil {
		if client, ok := value.(*kubernetes.Clientset); ok {
			return client
		}
	}
	return nil
}

func ConfigFromContext(ctx context.Context) *rest.Config {
	value := ctx.Value(configContextKey)
	if value != nil {
		if config, ok := value.(*rest.Config); ok {
			return config
		}
	}
	return nil
}

func RawConfigFromContext(ctx context.Context) *api.Config {
	value := ctx.Value(rawConfigContextKey)
	if value != nil {
		if config, ok := value.(*api.Config); ok {
			return config
		}
	}
	return nil
}

func RepoClientsetFromContext(ctx context.Context) repoclient.RepoClientset {
	value := ctx.Value(repoClientsetContextKey)
	if value != nil {
		if config, ok := value.(repoclient.RepoClientset); ok {
			return config
		}
	}
	return nil
}
