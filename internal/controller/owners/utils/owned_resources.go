package utils

import (
	"errors"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RefersToSameResource(a, b packagesv1alpha1.OwnedResourceRef) bool {
	return a.Group == b.Group &&
		a.Version == b.Version &&
		a.Kind == b.Kind &&
		a.Name == b.Name &&
		a.Namespace == b.Namespace
}

func Add(refs *[]packagesv1alpha1.OwnedResourceRef, newRefs ...packagesv1alpha1.OwnedResourceRef) bool {
	var changed bool
outer:
	for _, newRef := range newRefs {
		for _, ref := range *refs {
			if RefersToSameResource(ref, newRef) {
				continue outer
			}
		}
		*refs = append(*refs, newRef)
		changed = true
	}
	return changed
}

func Remove(refs *[]packagesv1alpha1.OwnedResourceRef, toRemove packagesv1alpha1.OwnedResourceRef) bool {
	for i, ref := range *refs {
		if RefersToSameResource(ref, toRemove) {
			*refs = append((*refs)[0:i], (*refs)[i+1:]...)
			return true
		}
	}
	return false
}

func MarkForDeletion(refs *[]packagesv1alpha1.OwnedResourceRef, toRemove packagesv1alpha1.OwnedResourceRef) bool {
	for i, ref := range *refs {
		if RefersToSameResource(ref, toRemove) {
			(*refs)[i].MarkedForDeletion = true
			return true
		}
	}
	return false
}

func AddOwnedResourceRef(
	scheme *runtime.Scheme,
	refs *[]packagesv1alpha1.OwnedResourceRef,
	obj client.Object,
) (bool, error) {
	if ref, err := ToOwnedResourceRef(scheme, obj); err != nil {
		return false, err
	} else {
		return Add(refs, ref), nil
	}
}

func RemoveOwnedResourceRef(refs *[]packagesv1alpha1.OwnedResourceRef, ref packagesv1alpha1.OwnedResourceRef) {
	for i, r := range *refs {
		if RefersToSameResource(r, ref) {
			*refs = append((*refs)[:i], (*refs)[i+1:]...)
			return
		}
	}
}

func ToOwnedResourceRef(scheme *runtime.Scheme, obj client.Object) (packagesv1alpha1.OwnedResourceRef, error) {
	if gvk, err := GetGVK(scheme, obj); err != nil {
		return packagesv1alpha1.OwnedResourceRef{}, err
	} else {
		return packagesv1alpha1.OwnedResourceRef{
			GroupVersionKind: gvk,
			Name:             obj.GetName(),
			Namespace:        obj.GetNamespace(),
		}, nil
	}
}

func OwnedResourceRefToObject(ref packagesv1alpha1.OwnedResourceRef) client.Object {
	obj := unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind(ref.GroupVersionKind))
	obj.SetName(ref.Name)
	obj.SetNamespace(ref.Namespace)
	return &obj
}

func GetGVK(scheme *runtime.Scheme, obj client.Object) (metav1.GroupVersionKind, error) {
	if gvks, _, err := scheme.ObjectKinds(obj); err != nil {
		return metav1.GroupVersionKind{}, err
	} else {
		for _, gvk := range gvks {
			if len(gvk.Kind) > 0 && len(gvk.Version) > 0 {
				return metav1.GroupVersionKind(gvk), nil
			}
		}
	}
	return metav1.GroupVersionKind{}, errors.New("could not find a usable GroupVersionKind for object")
}
