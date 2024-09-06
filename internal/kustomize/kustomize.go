package kustomize

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kustomize/api/krusty"
	kstypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/yaml"
)

const (
	kustomizationFileName = "kustomization.yaml"
	objectsFileName       = "resources.yaml"
)

// KustomizeObjects runs "kustomize" over a [client.Object] slice, given a partial [kstypes.Kustomization].
// This is made possible by using an in-memory filesystem, which means that a large number of resources might lead to
// spikes in memory usage.
// The Kustomization may not specify any resources or generators, but thinks like namePrefix, namespace, labels should
// all work.
func KustomizeObjects(kustomization kstypes.Kustomization, objects []client.Object) ([]client.Object, error) {
	if fs, err := createVirtualFilesys(kustomization, objects); err != nil {
		return nil, err
	} else if resMap, err := krusty.MakeKustomizer(krusty.MakeDefaultOptions()).Run(fs, "."); err != nil {
		return nil, err
	} else {
		resources := resMap.Resources()
		result := make([]client.Object, len(resources))
		for i, res := range resMap.Resources() {
			if data, err := res.MarshalJSON(); err != nil {
				return nil, err
			} else {
				result[i] = &unstructured.Unstructured{}
				if err := json.Unmarshal(data, result[i]); err != nil {
					return nil, err
				}
			}
		}
		return result, nil
	}
}

func createVirtualFilesys(kustomization kstypes.Kustomization, objs []client.Object) (filesys.FileSystem, error) {
	fs := filesys.MakeFsInMemory()

	if f, err := fs.Create(kustomizationFileName); err != nil {
		return nil, err
	} else {
		defer func() { _ = f.Close() }()
		kustomization.Resources = append(kustomization.Resources, objectsFileName)
		if data, err := yaml.Marshal(kustomization); err != nil {
			return nil, err
		} else if _, err := f.Write(data); err != nil {
			return nil, err
		}
	}

	if f, err := fs.Create(objectsFileName); err != nil {
		return nil, err
	} else {
		defer func() { _ = f.Close() }()
		for _, obj := range objs {
			if data, err := yaml.Marshal(obj); err != nil {
				return nil, err
			} else if _, err = fmt.Fprintln(f, "\n---"); err != nil {
				return nil, err
			} else if _, err = f.Write(data); err != nil {
				return nil, err
			}
		}
	}
	return fs, nil
}
