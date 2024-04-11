package telemetry

import (
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/posthog/posthog-go"
)

type OperatorTelemetry struct {
	disabled  bool
	clusterId string
	posthog   posthog.Client
	start     time.Time
	packages  map[string]time.Time // TODO is this multithreaded?
}

func ForOperator() *OperatorTelemetry {
	if false {
		// TODO if telemetry
		// TODO if is dev?
		return &OperatorTelemetry{disabled: true}
	}
	if ph, err := posthog.NewWithConfig(
		apiKey,
		posthog.Config{
			Endpoint: "https://eu.posthog.com",
		},
	); err != nil {
		// TODO ?
		return nil
	} else {
		return &OperatorTelemetry{
			posthog:   ph,
			clusterId: "TODO get cluster ID", // TODO
			start:     time.Now(),
			packages:  make(map[string]time.Time),
		}
	}
}

func (t *OperatorTelemetry) ReconcilePackage(pkg *v1alpha1.Package) {
	if t.disabled {
		// TODO this should be get live every time, because it could have been enabled/disabled in the meanwhile
		return
	}
	if lastReported, ok := t.packages[pkg.Name]; ok && lastReported.Add(time.Minute*5).After(time.Now()) { // TODO change to x hours
		return
	}
	// TODO important to track for package: name, version, installation time creation timestamp der CR (set_once), status
	// TODO important to track for cluster: operator version, kubernetes version, etc etc (see document)
	err := t.posthog.Enqueue(posthog.Capture{
		DistinctId: t.clusterId,
		Event:      "reconcile_package",
		Properties: map[string]any{
			// event properties:
			"package":                 pkg.Name,
			"desired_package_version": pkg.Spec.PackageInfo.Version,
			"actual_package_version":  pkg.Status.Version,
			// user properties:
			"$set": map[string]any{
				"version": "TODO", // TODO
			},
			"$set_once": map[string]any{
				"type":            "operator",
				"initial_version": "TODO", // TODO
			},
		},
	})
	if err == nil {
		t.packages[pkg.Name] = time.Now()
	}
}

func (t *OperatorTelemetry) ReportDelete(pkg *v1alpha1.Package) {
	// TODO if disabled blabla
	// TODO report only once every few minutes
	_ = t.posthog.Enqueue(posthog.Capture{
		DistinctId: t.clusterId,
		Event:      "delete_package",
		Properties: map[string]any{
			// event properties:
			"package":                 pkg.Name,
			"desired_package_version": pkg.Spec.PackageInfo.Version,
			"actual_package_version":  pkg.Status.Version,
			// user properties:
			"$set": map[string]any{
				"version": "TODO", // TODO
			},
			"$set_once": map[string]any{
				"type":            "operator",
				"initial_version": "TODO", // TODO
			},
		},
	})
}
