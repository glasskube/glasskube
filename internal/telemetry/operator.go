package telemetry

import (
	"context"
	"sync"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/telemetry/properties"
	"github.com/glasskube/glasskube/pkg/condition"
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

type managerRepositoryLister struct{ client.Client }

func (l *managerRepositoryLister) ListRepositories(ctx context.Context) (*v1alpha1.PackageRepositoryList, error) {
	var ls v1alpha1.PackageRepositoryList
	return &ls, l.List(ctx, &ls)
}

var operatorInstance *OperatorTelemetry

func InitWithManager(mgr manager.Manager) {
	operatorInstance = ForControllerManager(mgr)
}

func ForOperator() *OperatorTelemetry {
	return operatorInstance
}

func ForControllerManager(mgr manager.Manager) *OperatorTelemetry {
	t := OperatorTelemetry{
		startTimestamp:            time.Now(),
		packageReportTimes:        make(map[string]time.Time),
		packageReportMuteDuration: 8 * time.Hour,
		PropertyGetter: properties.PropertyGetter{
			NamespaceGetter:  &managerNamespaceGetter{mgr.GetClient()},
			NodeLister:       &managerNodeLister{mgr.GetClient()},
			RepositoryLister: &managerRepositoryLister{mgr.GetClient()},
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

func (t *OperatorTelemetry) ReportStart() {
	if !t.Enabled() || t.posthog == nil {
		return
	}
	_ = t.posthog.Enqueue(posthog.Capture{
		DistinctId: t.ClusterId(),
		Event:      "operator_started",
		Properties: properties.BuildProperties(
			properties.ForOperatorUser(t.PropertyGetter),
		),
	})
}

func (t *OperatorTelemetry) OnEvent(obj client.Object, status condition.Type, reason condition.Reason) {
	_ = t.posthog.Enqueue(posthog.Capture{
		DistinctId: t.ClusterId(),
		Event:      "status_conditions_changed",
		Properties: properties.BuildProperties(
			properties.ForOperatorUser(t.PropertyGetter),
			properties.FromMap(map[string]any{
				"obj_kind":          obj.GetObjectKind().GroupVersionKind().Kind,
				"obj_name":          obj.GetName(),
				"obj_status_type":   status,
				"obj_status_reason": reason,
			}),
		),
	})
}

func (t *OperatorTelemetry) ReconcilePackage(pkg ctrlpkg.Package) {
	if !t.Enabled() || t.posthog == nil {
		return
	}
	go func() {
		t.packageReportTimesMutex.Lock()
		defer t.packageReportTimesMutex.Unlock()
		if lastReported, ok := t.packageReportTimes[pkg.GetName()]; ok &&
			lastReported.Add(t.packageReportMuteDuration).After(time.Now()) {
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
			t.packageReportTimes[pkg.GetName()] = time.Now()
		}
	}()
}

func (t *OperatorTelemetry) ReportDelete(pkg ctrlpkg.Package) {
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
