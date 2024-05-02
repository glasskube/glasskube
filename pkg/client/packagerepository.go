package client

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var packageRepositoryGVR = v1alpha1.GroupVersion.WithResource("packagerepositories")

type PackageRepositoryInterface interface {
	Create(ctx context.Context, obj *v1alpha1.PackageRepository) error
	Update(ctx context.Context, obj *v1alpha1.PackageRepository) error
	Get(ctx context.Context, name string, obj *v1alpha1.PackageRepository) error
	GetAll(ctx context.Context, obj *v1alpha1.PackageRepositoryList) error
	Watch(ctx context.Context) (watch.Interface, error)
	Delete(ctx context.Context, obj *v1alpha1.PackageRepository, options metav1.DeleteOptions) error
}

type packageRepositoryClient struct {
	restClient rest.Interface
}

// Create implements PackageRepositoryInterface.
func (c *packageRepositoryClient) Create(ctx context.Context, obj *v1alpha1.PackageRepository) error {
	return c.restClient.Post().
		Resource(packageRepositoryGVR.Resource).
		Body(obj).Do(ctx).Into(obj)
}

// Update implements PackageRepositoryInterface.
func (c *packageRepositoryClient) Update(ctx context.Context, obj *v1alpha1.PackageRepository) error {
	return c.restClient.Put().
		Resource(packageRepositoryGVR.Resource).
		Name(obj.GetName()).
		Body(obj).
		Do(ctx).
		Into(obj)
}

// Watch implements PackageRepositoryInterface.
func (c *packageRepositoryClient) Watch(ctx context.Context) (watch.Interface, error) {
	opts := metav1.ListOptions{Watch: true}
	return c.restClient.Get().
		Resource(packageRepositoryGVR.Resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

// Get implements PackageRepositoryInterface.
func (c *packageRepositoryClient) Get(ctx context.Context, pkgName string, obj *v1alpha1.PackageRepository) error {
	return c.restClient.Get().
		Resource(packageRepositoryGVR.Resource).
		Name(pkgName).
		Do(ctx).Into(obj)
}

// GetAll implements PackageRepositoryInterface.
func (c *packageRepositoryClient) GetAll(ctx context.Context, result *v1alpha1.PackageRepositoryList) error {
	return c.restClient.Get().
		Resource(packageRepositoryGVR.Resource).
		Do(ctx).Into(result)
}

// Delete implements PackageRepositoryInterface.
func (c *packageRepositoryClient) Delete(
	ctx context.Context, obj *v1alpha1.PackageRepository, options metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Resource(packageRepositoryGVR.Resource).
		Name(obj.Name).
		Body(&options).
		Do(ctx).Into(nil)
}
