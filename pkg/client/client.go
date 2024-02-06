package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type PackageV1Alpha1Client struct {
	restClient rest.Interface
}

func (c *PackageV1Alpha1Client) Packages() PackageInterface {
	return &packageClient{
		restClient: c.restClient,
	}
}

func (c *PackageV1Alpha1Client) PackageInfos() PackageInfoInterface {
	return &packageInfoClient{restClient: c.restClient}
}

func New(cfg *rest.Config) (*PackageV1Alpha1Client, error) {
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
