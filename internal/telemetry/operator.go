package telemetry

import (
	"context"
	"sync"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/telemetry/properties"
	"github.com/posthog/posthog-go"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type OperatorTelemetry struct {
	properties.PropertyGetter
	posthog                   posthog.Client
	startTimestamp            time.Time
	packageReportTimes        map[string]time.Time
	packageReportTimesMutex   sync.Mutex
	packageReportMuteDuration time.Duration
}

type managerNamespaceGetter struct{ client.Client }

func (g *managerNamespaceGetter) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	var ns corev1.Namespace
	return &ns, g.Get(ctx, types.NamespacedName{Name: name}, &ns)
}

type managerNodeLister struct{ client.Client }

func (l *managerNodeLister) ListNodes(ctx context.Context) (*corev1.NodeList, error) {
	var nl corev1.NodeList
	return &nl, l.List(ctx, &nl)
}

func ForControllerManager(mgr manager.Manager) *OperatorTelemetry {
	t := OperatorTelemetry{
		startTimestamp:            time.Now(),
		packageReportTimes:        make(map[string]time.Time),
		packageReportMuteDuration: 5 * time.Second, // TODO change to x hours
		PropertyGetter: properties.PropertyGetter{
			NamespaceGetter: &managerNamespaceGetter{mgr.GetClient()},
			NodeLister:      &managerNodeLister{mgr.GetClient()},
		},
	}
	if discoveryClient, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig()); err == nil {
		t.PropertyGetter.DiscoveryClient = discoveryClient
	}
	if ph, err := posthog.NewWithConfig(apiKey, posthog.Config{Endpoint: endpoint}); err == nil {
		t.posthog = ph
	}
	return &t
}

func (t *OperatorTelemetry) ReconcilePackage(pkg *v1alpha1.Package) {
	if !t.Enabled() || t.posthog == nil {
		return
	}
	go func() {
		t.packageReportTimesMutex.Lock()
		defer t.packageReportTimesMutex.Unlock()
		if lastReported, ok := t.packageReportTimes[pkg.Name]; ok && lastReported.Add(t.packageReportMuteDuration).After(time.Now()) {
			return
		}
		err := t.posthog.Enqueue(posthog.Capture{
			DistinctId: t.ClusterId(),
			Event:      "reconcile_package",
			Properties: properties.BuildProperties(
				properties.ForOperatorUser(t.PropertyGetter),
				properties.FromPackage(pkg),
			),
		})
		if err == nil {
			t.packageReportTimes[pkg.Name] = time.Now()
		}
	}()
}

func (t *OperatorTelemetry) ReportDelete(pkg *v1alpha1.Package) {
	if !t.Enabled() || t.posthog == nil {
		return
	}
	_ = t.posthog.Enqueue(posthog.Capture{
		DistinctId: t.ClusterId(),
		Event:      "delete_package",
		Properties: properties.BuildProperties(
			properties.ForOperatorUser(t.PropertyGetter),
			properties.FromPackage(pkg),
		),
	})
}
