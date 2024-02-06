package client

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/rest"
)

var packageInfoGVR = v1alpha1.GroupVersion.WithResource("packageinfos")

type PackageInfoInterface interface {
	Get(ctx context.Context, name string, packageInfo *v1alpha1.PackageInfo) error
}

type packageInfoClient struct {
	restClient rest.Interface
}

// Get implements PackageInfoInterface.
func (c *packageInfoClient) Get(ctx context.Context, name string, packageInfo *v1alpha1.PackageInfo) error {
	return c.restClient.Get().
		Resource(packageInfoGVR.Resource).
		Name(name).
		Do(ctx).
		Into(packageInfo)
}
