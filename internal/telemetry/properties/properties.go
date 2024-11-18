package properties

import (
	"context"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"

	"github.com/glasskube/glasskube/internal/telemetry/annotations"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
)

type NodeLister interface {
	ListNodes(ctx context.Context) (*corev1.NodeList, error)
}

type NamespaceGetter interface {
	GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error)
}

type RepositoryLister interface {
	ListRepositories(ctx context.Context) (*v1alpha1.PackageRepositoryList, error)
}

type PropertyGetter struct {
	NodeLister       NodeLister
	NamespaceGetter  NamespaceGetter
	DiscoveryClient  discovery.DiscoveryInterface
	RepositoryLister RepositoryLister
}

type ClusterProperties struct {
	kubernetesVersion string
	provider          string
	nnodes            int
	gitopsMode        bool
}

type RepositoryProperties struct {
	nrepositories       int
	nrepositoriesAuth   int
	customRepoAsDefault bool
}

func (g PropertyGetter) Enabled() bool {
	if g.NamespaceGetter != nil {
		if ns, err := g.NamespaceGetter.GetNamespace(context.Background(), "glasskube-system"); err == nil {
			return annotations.IsTelemetryEnabled(ns.Annotations)
		}
	}
	return false
}

func (g PropertyGetter) ClusterId() string {
	if g.NamespaceGetter != nil {
		if ns, err := g.NamespaceGetter.GetNamespace(context.Background(), "glasskube-system"); err == nil {
			return ns.Annotations[annotations.TelemetryIdAnnotation]
		}
	}
	return ""
}

func (g PropertyGetter) ClusterProperties() (p ClusterProperties) {
	if g.DiscoveryClient != nil {
		if versionInfo, err := g.DiscoveryClient.ServerVersion(); err == nil {
			p.kubernetesVersion = versionInfo.GitVersion
		}
	}
	if g.NodeLister != nil {
		if nodes, err := g.NodeLister.ListNodes(context.Background()); err == nil {
			p.nnodes = len(nodes.Items)
			for _, node := range nodes.Items {
				// ProviderID is the ID assigend to the node by the provider.
				// It usually has the format <provider>://<nodeId>.
				splits := strings.SplitN(node.Spec.ProviderID, "://", 2)
				if len(splits) > 1 {
					p.provider = splits[0]
					break
				}
			}
		}
	}
	if g.NamespaceGetter != nil {
		if ns, err := g.NamespaceGetter.GetNamespace(context.Background(), "glasskube-system"); err == nil {
			p.gitopsMode = annotations.IsGitopsModeEnabled(ns.Annotations)
		}
	}
	return
}

func (g PropertyGetter) RepositoryProperties() (p RepositoryProperties) {
	if g.RepositoryLister != nil {
		if ls, err := g.RepositoryLister.ListRepositories(context.Background()); err == nil {
			p.nrepositories = len(ls.Items)
			for _, repo := range ls.Items {
				if repo.Spec.Auth != nil {
					p.nrepositoriesAuth = p.nrepositoriesAuth + 1
				}
				if repo.IsDefaultRepository() && !repo.IsGlasskubeRepo() {
					p.customRepoAsDefault = true
				}
			}
		}
	}
	return
}
