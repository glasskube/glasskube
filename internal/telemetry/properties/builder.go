package properties

import (
	"runtime"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/posthog/posthog-go"
	"k8s.io/apimachinery/pkg/api/meta"
)

type PropertiesBuilderFn func(p posthog.Properties) posthog.Properties

func BuildProperties(fns ...PropertiesBuilderFn) posthog.Properties {
	properties := posthog.NewProperties()
	for _, fn := range fns {
		properties = fn(properties)
	}
	return properties
}

func ForClientUser(pg PropertyGetter) PropertiesBuilderFn {
	return func(p posthog.Properties) posthog.Properties {
		cp := pg.ClusterProperties()
		return p.
			Set("$set", map[string]any{
				"version": config.Version,
			}).
			Set("$set_once", map[string]any{
				"type":            "client",
				"initial_version": config.Version,
				"os":              runtime.GOOS,
				"architecture":    runtime.GOARCH,
			}).
			Set("cluster_id", pg.ClusterId()).
			Set("cluster_k8s_version", cp.kubernetesVersion).
			Set("cluster_provider", cp.provider).
			Set("cluster_nnodes", cp.nnodes)
	}
}

func ForOperatorUser(pg PropertyGetter) PropertiesBuilderFn {
	// TODO: see document for additional required properties
	return func(p posthog.Properties) posthog.Properties {
		cp := pg.ClusterProperties()
		return p.
			Set("$set", map[string]any{
				"version":     config.Version,
				"k8s_version": cp.kubernetesVersion,
				"provider":    cp.provider,
				"nnodes":      cp.nnodes,
			}).
			Set("$set_once", map[string]any{
				"type":            "operator",
				"initial_version": config.Version,
			})
	}
}

func FromMap(data map[string]any) PropertiesBuilderFn {
	return func(p posthog.Properties) posthog.Properties {
		for k, v := range data {
			p.Set(k, v)
		}
		return p
	}
}

func FromPackage(pkg *v1alpha1.Package) PropertiesBuilderFn {
	return func(p posthog.Properties) posthog.Properties {
		p.Set("package_name", pkg.Spec.PackageInfo.Name).
			Set("package_version_desired", pkg.Spec.PackageInfo.Version).
			Set("package_version_actual", pkg.Status.Version).
			// TODO: set_once ?
			Set("package_creation_timestamp", pkg.CreationTimestamp)
		if c := meta.FindStatusCondition(pkg.Status.Conditions, string(condition.Ready)); c != nil {
			p.Set("package_ready_status", c.Status)
			p.Set("package_ready_reason", c.Reason)
		}
		return p
	}
}
