package properties

import (
	"runtime"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
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

func ForClientUser(pg PropertyGetter, includeCluster bool) PropertiesBuilderFn {
	return func(p posthog.Properties) posthog.Properties {
		var cp ClusterProperties
		if includeCluster {
			cp = pg.ClusterProperties()
			p.Set("cluster_id", pg.ClusterId()).
				Set("cluster_k8s_version", cp.kubernetesVersion).
				Set("cluster_provider", cp.provider).
				Set("cluster_nnodes", cp.nnodes)
		}
		p.Set("$set", map[string]any{
			"version": config.Version,
		}).Set("$set_once", map[string]any{
			"type":            "client",
			"initial_version": config.Version,
			"os":              runtime.GOOS,
			"architecture":    runtime.GOARCH,
		}).Set("version", config.Version)
		return p
	}
}

func ForOperatorUser(pg PropertyGetter) PropertiesBuilderFn {
	return func(p posthog.Properties) posthog.Properties {
		cp := pg.ClusterProperties()
		rp := pg.RepositoryProperties()
		return p.
			Set("$set", map[string]any{
				"version":             config.Version,
				"k8s_version":         cp.kubernetesVersion,
				"provider":            cp.provider,
				"nnodes":              cp.nnodes,
				"nrepositories":       rp.nrepositories,
				"nrepositories_auth":  rp.nrepositoriesAuth,
				"custom_default_repo": rp.customRepoAsDefault,
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

func FromPackage(pkg ctrlpkg.Package) PropertiesBuilderFn {
	return func(p posthog.Properties) posthog.Properties {
		p.Set("package_kind", pkg.GroupVersionKind().Kind).
			Set("package_name", pkg.GetSpec().PackageInfo.Name).
			Set("package_version_desired", pkg.GetSpec().PackageInfo.Version).
			Set("package_version_actual", pkg.GetStatus().Version).
			// TODO: set_once ?
			Set("package_creation_timestamp", pkg.GetCreationTimestamp()).
			Set("package_auto_update", pkg.AutoUpdatesEnabled()).
			Set("package_repository_name", pkg.GetSpec().PackageInfo.RepositoryName)
		if c := meta.FindStatusCondition(pkg.GetStatus().Conditions, string(condition.Ready)); c != nil {
			p.Set("package_ready_status", c.Status)
			p.Set("package_ready_reason", c.Reason)
		}
		return p
	}
}
