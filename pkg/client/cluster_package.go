//nolint:dupl
package client

import (
	"context"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var clusterPackageGVR = v1alpha1.GroupVersion.WithResource("clusterpackages")

type clusterPackageClient struct {
	restClient rest.Interface
}

func (c *clusterPackageClient) Create(
	ctx context.Context, pkg *v1alpha1.ClusterPackage, opts metav1.CreateOptions) error {
	return c.restClient.Post().
		Resource(clusterPackageGVR.Resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(pkg).Do(ctx).Into(pkg)
}

// Update implements PackageInterface.
func (c *clusterPackageClient) Update(ctx context.Context, p *v1alpha1.ClusterPackage) error {
	return c.restClient.Put().
		Resource(clusterPackageGVR.Resource).
		Name(p.GetName()).
		Body(p).
		Do(ctx).
		Into(p)
}

func (c *clusterPackageClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.restClient.Get().
		Resource(clusterPackageGVR.Resource).
		Timeout(timeout).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *clusterPackageClient) Get(ctx context.Context, pkgName string, result *v1alpha1.ClusterPackage) error {
	return c.restClient.Get().
		Resource(clusterPackageGVR.Resource).
		Name(pkgName).
		Do(ctx).Into(result)
}

func (c *clusterPackageClient) GetAll(ctx context.Context, result *v1alpha1.ClusterPackageList) error {
	return c.restClient.Get().
		Resource(clusterPackageGVR.Resource).
		Do(ctx).Into(result)
}

func (c *clusterPackageClient) Delete(
	ctx context.Context, pkg *v1alpha1.ClusterPackage, options metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Resource(clusterPackageGVR.Resource).
		Name(pkg.Name).
		Body(&options).
		Do(ctx).Into(nil)
}
