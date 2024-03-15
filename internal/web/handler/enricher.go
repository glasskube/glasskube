package handler

import (
	"net/http"

	"github.com/glasskube/glasskube/pkg/client"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

type ContextDataSupplier interface {
	RestConfig() *rest.Config
	RawConfig() *api.Config
	Client() client.PackageV1Alpha1Client
}

type ContextEnrichingHandler struct {
	Source  ContextDataSupplier
	Handler http.Handler
}

func (enricher *ContextEnrichingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := client.SetupContextWithClient(r.Context(),
		enricher.Source.RestConfig(),
		enricher.Source.RawConfig(),
		enricher.Source.Client())
	enricher.Handler.ServeHTTP(w, r.WithContext(ctx))
}
