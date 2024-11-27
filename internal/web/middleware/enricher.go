package middleware

import (
	"net/http"

	"github.com/glasskube/glasskube/internal/web/types"

	"github.com/glasskube/glasskube/internal/clicontext"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	webcontext "github.com/glasskube/glasskube/internal/web/context"
	"github.com/glasskube/glasskube/pkg/client"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

type ContextDataSupplier interface {
	RestConfig() *rest.Config
	RawConfig() *api.Config
	Client() client.PackageV1Alpha1Client
	K8sClient() *kubernetes.Clientset
	RepoClient() repoclient.RepoClientset
	CoreListers() *types.CoreListers
}

type ContextEnrichingHandler struct {
	Source  ContextDataSupplier
	Handler http.Handler
}

func (enricher *ContextEnrichingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := clicontext.SetupContextWithClient(r.Context(),
		enricher.Source.RestConfig(),
		enricher.Source.RawConfig(),
		enricher.Source.Client(),
		enricher.Source.K8sClient())
	ctx = clicontext.ContextWithRepositoryClientset(ctx, enricher.Source.RepoClient())
	ctx = webcontext.ContextWithCoreListers(ctx, enricher.Source.CoreListers())
	enricher.Handler.ServeHTTP(w, r.WithContext(ctx))
}
