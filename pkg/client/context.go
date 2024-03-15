package client

import (
	"context"

	"k8s.io/client-go/tools/clientcmd/api"

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
	if pkgClient, err := New(config); err != nil {
		return nil, err
	} else {
		return SetupContextWithClient(ctx, config, rawConfig, pkgClient), nil
	}
}

func SetupContextWithClient(
	ctx context.Context,
	config *rest.Config,
	rawConfig *api.Config,
	client PackageV1Alpha1Client,
) context.Context {
	ctx = context.WithValue(ctx, clientContextKey, client)
	ctx = context.WithValue(ctx, configContextKey, config)
	ctx = context.WithValue(ctx, rawConfigContextKey, rawConfig)
	return ctx
}

func FromContext(ctx context.Context) PackageV1Alpha1Client {
	value := ctx.Value(clientContextKey)
	if value != nil {
		if client, ok := value.(PackageV1Alpha1Client); ok {
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
