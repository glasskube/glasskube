package cliutils

import (
	"context"

	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/pkg/client"
	"k8s.io/client-go/kubernetes"
)

func PackageClient(ctx context.Context) client.PackageV1Alpha1Client {
	return clicontext.PackageClientFromContext(ctx)
}

func KubernetesClient(ctx context.Context) *kubernetes.Clientset {
	return clicontext.KubernetesClientFromContext(ctx)
}

func DependencyManager(ctx context.Context) *dependency.DependendcyManager {
	return dependency.NewDependencyManager(
		clientadapter.NewPackageClientAdapter(PackageClient(ctx)),
		RepositoryClientset(ctx),
	)
}

func ValueResolver(ctx context.Context) *manifestvalues.Resolver {
	return manifestvalues.NewResolver(
		clientadapter.NewPackageClientAdapter(PackageClient(ctx)),
		clientadapter.NewKubernetesClientAdapter(KubernetesClient(ctx)),
	)
}

func RepositoryClientset(ctx context.Context) repoclient.RepoClientset {
	if existing := clicontext.RepoClientsetFromContext(ctx); existing != nil {
		return existing
	}
	return repoclient.NewClientset(
		clientadapter.NewPackageClientAdapter(PackageClient(ctx)),
		clientadapter.NewKubernetesClientAdapter(KubernetesClient(ctx)),
	)
}
