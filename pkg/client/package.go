package client

import (
	"context"
	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type PackageV1Alpha1Client struct {
	restClient rest.Interface
}

type PackageInterface interface {
	Create(ctx context.Context, p *v1alpha1.Package) error
	Get(ctx context.Context, pkgName string, p *v1alpha1.Package) error
	GetAll(ctx context.Context, result *v1alpha1.PackageList) error
	Watch(ctx context.Context) (watch.Interface, error)
	Delete(ctx context.Context, pkg *v1alpha1.Package) error
}

type packageClient struct {
	restClient rest.Interface
}

func NewPackageClient(cfg *rest.Config) (*PackageV1Alpha1Client, error) {
	pkgRestConfig := *cfg
	pkgRestConfig.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	pkgRestConfig.APIPath = "/apis"
	pkgRestConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	restClient, err := rest.RESTClientFor(&pkgRestConfig)
	if err != nil {
		return nil, err
	}
	return &PackageV1Alpha1Client{restClient: restClient}, err
}

func (c *PackageV1Alpha1Client) Packages() PackageInterface {
	return &packageClient{
		restClient: c.restClient,
	}
}

func (c *packageClient) Create(ctx context.Context, pkg *v1alpha1.Package) error {
	return c.restClient.Post().
		Resource(PackageGVR.Resource).
		Body(pkg).Do(ctx).Into(pkg)
}

func (c *packageClient) Watch(ctx context.Context) (watch.Interface, error) {
	opts := metav1.ListOptions{Watch: true}
	return c.restClient.Get().
		Resource(PackageGVR.Resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *packageClient) Get(ctx context.Context, pkgName string, result *v1alpha1.Package) error {
	return c.restClient.Get().
		Resource(PackageGVR.Resource).
		Name(pkgName).
		Do(ctx).Into(result)
}

func (c *packageClient) GetAll(ctx context.Context, result *v1alpha1.PackageList) error {
	return c.restClient.Get().
		Resource(PackageGVR.Resource).
		Do(ctx).Into(result)
}

func (c *packageClient) Delete(ctx context.Context, pkg *v1alpha1.Package) error {
	return c.restClient.Delete().
		Resource(PackageGVR.Resource).
		Name(pkg.Name).
		Do(ctx).Into(nil)
}

// NewPackage instantiates a new v1alpha1.Package struct with the given package name
func NewPackage(packageName string) *v1alpha1.Package {
	return &v1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: packageName,
		},
		Spec: v1alpha1.PackageSpec{
			PackageInfo: v1alpha1.PackageInfoTemplate{
				Name: packageName,
			},
		},
	}
}
