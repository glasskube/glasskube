package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type PackageV1Alpha1Client interface {
	Packages() PackageInterface
	PackageInfos() PackageInfoInterface
	PackageRepositories() PackageRepositoryInterface
	WithStores(packageStore cache.Store, packageInfoStore cache.Store) PackageV1Alpha1Client
}

type baseClient struct {
	restClient rest.Interface
}

func (c *baseClient) Packages() PackageInterface {
	return &packageClient{
		restClient: c.restClient,
	}
}

func (c *baseClient) PackageInfos() PackageInfoInterface {
	return &packageInfoClient{restClient: c.restClient}
}

func (c *baseClient) PackageRepositories() PackageRepositoryInterface {
	return &packageRepositoryClient{restClient: c.restClient}
}

func (c *baseClient) WithStores(packageStore cache.Store, packageInfoStore cache.Store) PackageV1Alpha1Client {
	return &packageCacheClient{PackageV1Alpha1Client: c, packageStore: packageStore, packageInfoStore: packageInfoStore}
}

func New(cfg *rest.Config) (PackageV1Alpha1Client, error) {
	pkgRestConfig := *cfg
	pkgRestConfig.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	pkgRestConfig.APIPath = "/apis"
	pkgRestConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	restClient, err := rest.RESTClientFor(&pkgRestConfig)
	if err != nil {
		return nil, err
	}
	return &baseClient{restClient: restClient}, err
}

func NewOrDie(cfg *rest.Config) PackageV1Alpha1Client {
	if client, err := New(cfg); err != nil {
		panic(err)
	} else {
		return client
	}
}

func init() {
	_ = v1alpha1.AddToScheme(scheme.Scheme)
}
