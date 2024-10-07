package plain

import (
	"slices"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/kustomize"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kstypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

var wellKnownLabelFieldSpecs = kstypes.FsSlice{
	// Reference: https://www.keycloak.org/operator/advanced-configuration#_pod_template
	{
		Gvk:                resid.Gvk{Group: "k8s.keycloak.org", Version: "v2alpha1", Kind: "Keycloak"},
		Path:               "/spec/unsupported/podTemplate/metadata/labels",
		CreateIfNotPresent: true,
	},
}

func prefixAndUpdateReferences(
	pkg ctrlpkg.Package,
	manifest *v1alpha1.PackageManifest,
	objects []client.Object,
) ([]client.Object, error) {
	if pkg.IsNamespaceScoped() {
		kustomization := createKustomization(pkg)
		transitiveObjects, err := createTransitiveObjects(manifest.TransitiveResources)
		if err != nil {
			return nil, err
		}
		if result, err := kustomize.KustomizeObjects(kustomization, append(objects, transitiveObjects...)); err != nil {
			return nil, err
		} else {
			// strip transitive objects from the result and free them up for GC
			return slices.Clip(result[0 : len(result)-len(transitiveObjects)]), nil
		}
	} else {
		return objects, nil
	}
}

func createKustomization(pkg ctrlpkg.Package) kstypes.Kustomization {
	return kstypes.Kustomization{
		Namespace:  pkg.GetNamespace(),
		NamePrefix: pkg.GetName() + "-",
		Labels: []kstypes.Label{
			{
				Pairs: map[string]string{
					v1alpha1.LabelPackageName:         pkg.GetSpec().PackageInfo.Name,
					v1alpha1.LabelPackageInstanceName: pkg.GetName(),
				},
				FieldSpecs:       wellKnownLabelFieldSpecs,
				IncludeSelectors: true,
				IncludeTemplates: true,
			},
		},
	}
}

func createTransitiveObjects(resList []v1.TypedLocalObjectReference) ([]client.Object, error) {
	result := make([]client.Object, len(resList))
	for i, res := range resList {
		if obj, err := createTransitiveObject(res); err != nil {
			return nil, err
		} else {
			result[i] = obj
		}
	}
	return result, nil
}

// createTransitiveObject creates a "stub" [client.Object] for a given [v1alpha1.TransitiveResource].
// This Object can not be applied in a cluster and serve only as "reference", so the kustomize name reference
// transformer updates any references to them.
func createTransitiveObject(res v1.TypedLocalObjectReference) (client.Object, error) {
	var result unstructured.Unstructured
	if res.APIGroup == nil {
		result.SetGroupVersionKind(schema.GroupVersionKind{Kind: res.Kind})
	} else if gv, err := schema.ParseGroupVersion(*res.APIGroup); err != nil {
		return nil, err
	} else {
		result.SetGroupVersionKind(gv.WithKind(res.Kind))
	}
	result.SetName(res.Name)
	return &result, nil
}
