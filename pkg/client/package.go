package client

import (
	"context"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const packagesResource = "packages"

type packageClient struct {
	restClient rest.Interface
	ns         string
}

// Create implements PackageInterface.
func (p *packageClient) Create(ctx context.Context, target *v1alpha1.Package, opts v1.CreateOptions) error {
	return p.restClient.Post().
		Namespace(p.ns).
		Resource(packagesResource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(target).Do(ctx).Into(target)
}

// Delete implements PackageInterface.
func (p *packageClient) Delete(ctx context.Context, target *v1alpha1.Package, options v1.DeleteOptions) error {
	return p.restClient.Delete().
		Namespace(p.ns).
		Resource(packagesResource).
		Name(target.Name).
		Body(&options).
		Do(ctx).Into(nil)
}

// Get implements PackageInterface.
func (p *packageClient) Get(ctx context.Context, name string, target *v1alpha1.Package) error {
	return p.restClient.Get().
		Namespace(p.ns).
		Resource(packagesResource).
		Name(name).
		Do(ctx).Into(target)
}

// GetAll implements PackageInterface.
func (p *packageClient) GetAll(ctx context.Context, target *v1alpha1.PackageList) error {
	return p.restClient.Get().
		Namespace(p.ns).
		Resource(packagesResource).
		Do(ctx).Into(target)
}

// Update implements PackageInterface.
func (p *packageClient) Update(ctx context.Context, target *v1alpha1.Package) error {
	return p.restClient.Put().
		Namespace(p.ns).
		Resource(packagesResource).
		Name(target.Name).
		Body(target).
		Do(ctx).
		Into(target)
}

// Watch implements PackageInterface.
func (p *packageClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return p.restClient.Get().
		Namespace(p.ns).
		Resource(packagesResource).
		Timeout(timeout).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
