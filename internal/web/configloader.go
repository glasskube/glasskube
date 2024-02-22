package web

import (
	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

type configLoader interface {
	LoadConfig() (*rest.Config, *api.Config, error)
}

type defaultConfigLoader struct {
	Kubeconfig string
}

func (l *defaultConfigLoader) LoadConfig() (*rest.Config, *api.Config, error) {
	return kubeconfig.New(l.Kubeconfig)
}

type bytesConfigLoader struct {
	data []byte
}

func (l *bytesConfigLoader) LoadConfig() (*rest.Config, *api.Config, error) {
	return kubeconfig.FromBytes(l.data)
}
