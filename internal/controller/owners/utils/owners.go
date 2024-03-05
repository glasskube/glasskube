package utils

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ObjHasOwner(obj, owner client.Object) (bool, error) {
	refs := obj.GetOwnerReferences()
	for _, ref := range refs {
		if ref.Name != owner.GetName() {
			continue
		}
		if gv, err := schema.ParseGroupVersion(ref.APIVersion); err != nil {
			return false, err
		} else if owner.GetObjectKind().GroupVersionKind() == gv.WithKind(ref.Kind) {
			return true, nil
		}
	}
	return false, nil
}
