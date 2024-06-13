package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type baseClientset struct {
	restClient rest.Interface
}

// Packages implements PackageV1Alpha1Client.
func (c *baseClientset) Packages(namespace string) PackageInterface {
	return &packageClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *baseClientset) ClusterPackages() ClusterPackageInterface {
	return &clusterPackageClient{
		restClient: c.restClient,
	}
}

func (c *baseClientset) PackageInfos() PackageInfoInterface {
	return &packageInfoClient{restClient: c.restClient}
}

func (c *baseClientset) PackageRepositories() PackageRepositoryInterface {
	return &packageRepositoryClient{restClient: c.restClient}
}

func (c *baseClientset) WithStores(
	clusterPackageStore cache.Store,
	packageStore cache.Store,
	packageInfoStore cache.Store,
	packageRepositoryStore cache.Store,
) PackageV1Alpha1Client {
	return &cacheClientset{
		PackageV1Alpha1Client:  c,
		clusterPackageStore:    clusterPackageStore,
		packageStore:           packageStore,
		packageInfoStore:       packageInfoStore,
		packageRepositoryStore: packageRepositoryStore,
	}
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
	return &baseClientset{restClient: restClient}, err
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
