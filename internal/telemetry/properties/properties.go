package properties

import (
	"context"
	"strings"

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

type PropertyGetter struct {
	NodeLister      NodeLister
	NamespaceGetter NamespaceGetter
	DiscoveryClient discovery.DiscoveryInterface
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

func (g PropertyGetter) ClusterVersion() string {
	if g.DiscoveryClient != nil {
		if versionInfo, err := g.DiscoveryClient.ServerVersion(); err == nil {
			return versionInfo.GitVersion
		}
	}
	return ""
}

func (g PropertyGetter) ClusterProvider() string {
	if g.NodeLister != nil {
		if nodes, err := g.NodeLister.ListNodes(context.Background()); err == nil {
			for _, node := range nodes.Items {
				// ProviderID is the ID assigend to the node by the provider.
				// It usually has the format <provider>://<nodeId>.
				splits := strings.SplitN(node.Spec.ProviderID, "://", 2)
				if len(splits) > 1 {
					return splits[0]
				}
			}
		}
	}
	return ""
}
