package plain

import (
	"slices"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/kustomize"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kstypes "sigs.k8s.io/kustomize/api/types"
)

func prefixAndUpdateReferences(
	pkg ctrlpkg.Package,
	manifest *v1alpha1.PackageManifest,
	objects []client.Object,
) ([]client.Object, error) {
	if pkg.IsNamespaceScoped() {
		kustomization := createKustomization(pkg)
		transitiveObjects := createTransitiveObjects(manifest.TransitiveResources)
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
				IncludeSelectors: true,
				IncludeTemplates: true,
			},
		},
	}
}

func createTransitiveObjects(resList []v1alpha1.TransitiveResource) []client.Object {
	result := make([]client.Object, len(resList))
	for i, res := range resList {
		result[i] = createTransitiveObject(res)
	}
	return result
}

// createTransitiveObject creates a "stub" [client.Object] for a given [v1alpha1.TransitiveResource].
// This Object can not be applied in a cluster and serve only as "reference", so the kustomize name reference
// transformer updates any references to them.
func createTransitiveObject(res v1alpha1.TransitiveResource) client.Object {
	var result unstructured.Unstructured
	result.SetName(res.Name)
	result.SetGroupVersionKind(schema.GroupVersionKind(res.GroupVersionKind))
	return &result
}
