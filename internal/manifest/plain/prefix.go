package plain

import (
	"encoding/json"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kustomize/api/krusty"
	kstypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/yaml"
)

func createFS(
	pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest, objs []client.Object,
) (filesys.FileSystem, error) {
	fsys := filesys.MakeFsInMemory()

	if f, err := fsys.Create("kustomization.yaml"); err != nil {
		return nil, err
	} else {
		defer func() { _ = f.Close() }()
		kustomization := kstypes.Kustomization{
			Namespace:  pkg.GetNamespace(),
			NamePrefix: pkg.GetName() + "-",
			Labels: []kstypes.Label{
				{
					Pairs: map[string]string{
						"packages.glasskube.dev/package":  pkg.GetSpec().PackageInfo.Name,
						"packages.glasskube.dev/instance": pkg.GetName(),
					},
					IncludeSelectors: true,
					IncludeTemplates: true,
				},
			},
			Resources: []string{"resources.yaml"},
		}
		if data, err := yaml.Marshal(kustomization); err != nil {
			return nil, err
		} else if _, err := f.Write(data); err != nil {
			return nil, err
		}
	}

	if f, err := fsys.Create("resources.yaml"); err != nil {
		return nil, err
	} else {
		defer func() { _ = f.Close() }()
		for _, obj := range objs {
			if data, err := yaml.Marshal(obj); err != nil {
				return nil, err
			} else if _, err = fmt.Fprintln(f, "---"); err != nil {
				return nil, err
			} else if _, err = f.Write(data); err != nil {
				return nil, err
			}
		}

		for _, obj := range manifest.TransitiveResources {
			if data, err := yaml.Marshal(map[string]any{
				"apiVersion": metav1.GroupVersion{Group: obj.Group, Version: obj.Version}.String(),
				"kind":       obj.Kind,
				"metadata":   map[string]any{"name": obj.Name},
			}); err != nil {
				return nil, err
			} else if _, err = fmt.Fprintln(f, "---"); err != nil {
				return nil, err
			} else if _, err = f.Write(data); err != nil {
				return nil, err
			}
		}
	}
	return fsys, nil
}

func prefixAndUpdateReferences(
	pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest, objects []client.Object,
) ([]client.Object, error) {
	if pkg.IsNamespaceScoped() {
		if fsys, err := createFS(pkg, manifest, objects); err != nil {
			return nil, err
		} else if resMap, err := krusty.MakeKustomizer(krusty.MakeDefaultOptions()).Run(fsys, "."); err != nil {
			return nil, err
		} else {
			resources := resMap.Resources()
			result := make([]client.Object, len(resources))
			for i, res := range resMap.Resources() {
				if data, err := res.MarshalJSON(); err != nil {
					return nil, err
				} else {
					var u unstructured.Unstructured
					if err := json.Unmarshal(data, &u); err != nil {
						return nil, err
					}
					result[i] = &u
				}
			}
			return result[0 : len(result)-len(manifest.TransitiveResources)], nil
		}
	} else {
		return objects, nil
	}
}
