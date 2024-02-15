package client

import (
	"context"

	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type (
	contextKey int
)

const (
	clientContextKey contextKey = iota
	configContextKey
	rawConfigContextKey
)

func SetupContext(ctx context.Context, config *rest.Config, rawConfig *api.Config) (context.Context, error) {
	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}
	pkgClient, err := New(config)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientContextKey, pkgClient)
	ctx = context.WithValue(ctx, configContextKey, config)
	ctx = context.WithValue(ctx, rawConfigContextKey, rawConfig)
	return ctx, nil
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
