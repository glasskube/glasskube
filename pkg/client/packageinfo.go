package client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/rest"
)

var packageInfoGVR = v1alpha1.GroupVersion.WithResource("packageinfos")

type packageInfoClient struct {
	restClient rest.Interface
}

// GetAll implements PackageInfoInterface.
func (c *packageInfoClient) GetAll(ctx context.Context, result *v1alpha1.PackageInfoList) error {
	return c.restClient.Get().
		Resource(packageInfoGVR.Resource).
		Do(ctx).Into(result)
}

// Get implements PackageInfoInterface.
func (c *packageInfoClient) Get(ctx context.Context, name string, packageInfo *v1alpha1.PackageInfo) error {
	return c.restClient.Get().
		Resource(packageInfoGVR.Resource).
		Name(name).
		Do(ctx).
		Into(packageInfo)
}

func (c *packageInfoClient) Watch(ctx context.Context) (watch.Interface, error) {
	opts := metav1.ListOptions{Watch: true}
	return c.restClient.Get().
		Resource(packageInfoGVR.Resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
