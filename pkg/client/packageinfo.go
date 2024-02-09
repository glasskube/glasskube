package client

import (
	"context"

	packagev1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/rest"
)

var packageInfoGVR = packagev1alpha1.GroupVersion.WithResource("packageinfos")

type PackageInfoInterface interface {
	Get(ctx context.Context, name string, packageInfo *packagev1alpha1.PackageInfo) error
	GetAll(ctx context.Context, packageInfo *packagev1alpha1.PackageInfoList) error
}

type packageInfoClient struct {
	restClient rest.Interface
}

// GetAll implements PackageInfoInterface.
func (c *packageInfoClient) GetAll(ctx context.Context, result *packagev1alpha1.PackageInfoList) error {
	return c.restClient.Get().
		Resource(packageInfoGVR.Resource).
		Do(ctx).Into(result)
}

// Get implements PackageInfoInterface.
func (c *packageInfoClient) Get(ctx context.Context, name string, packageInfo *packagev1alpha1.PackageInfo) error {
	return c.restClient.Get().
		Resource(packageInfoGVR.Resource).
		Name(name).
		Do(ctx).
		Into(packageInfo)
}
