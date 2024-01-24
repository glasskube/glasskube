package client

import (
	"context"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type (
	contextKey int
)

const (
	clientContextKey contextKey = iota
)

var PackageGVR = v1alpha1.GroupVersion.WithResource("packages")

func SetupContext(ctx context.Context, config *rest.Config) (context.Context, error) {
	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}
	pkgClient, err := NewPackageClient(config)
	if err != nil {
		return nil, err
	}
	log.SetLogger(logr.New(log.NullLogSink{}))
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
