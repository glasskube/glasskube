package controller

import (
	"regexp"
	"strings"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	resourceNameRegex = regexp.MustCompile(`[^\w.-]`)
)

func objHasOwner(obj, owner client.Object) (bool, error) {
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

func escapeResourceName(name string) string {
	return strings.ToLower(resourceNameRegex.ReplaceAllString(name, "--"))
}

func refersToSameResource(a, b packagesv1alpha1.OwnedResourceRef) bool {
	return a.Group == b.Group &&
		a.Version == b.Version &&
		a.Kind == b.Kind &&
		a.Name == b.Name &&
		a.Namespace == b.Namespace
}

func addOwnedResourceRef(refs *[]packagesv1alpha1.OwnedResourceRef, kind schema.ObjectKind, obj metav1.Object) bool {
	newRef := toOwnedResourceRef(kind, obj)
	for _, ref := range *refs {
		if refersToSameResource(ref, newRef) {
			return false
		}
	}
	*refs = append(*refs, newRef)
	return true
}

func removeOwnedResourceRef(refs *[]packagesv1alpha1.OwnedResourceRef, ref packagesv1alpha1.OwnedResourceRef) {
	for i, r := range *refs {
		if refersToSameResource(r, ref) {
			*refs = append((*refs)[:i], (*refs)[i+1:]...)
			return
		}
	}
}

func toOwnedResourceRef(kind schema.ObjectKind, obj metav1.Object) packagesv1alpha1.OwnedResourceRef {
	return packagesv1alpha1.OwnedResourceRef{
		GroupVersionKind: metav1.GroupVersionKind(kind.GroupVersionKind()),
		Name:             obj.GetName(),
		Namespace:        obj.GetNamespace(),
	}
}

func ownedResourceRefToObject(ref packagesv1alpha1.OwnedResourceRef) client.Object {
	obj := unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind(ref.GroupVersionKind))
	obj.SetName(ref.Name)
	obj.SetNamespace(ref.Namespace)
	return &obj
}
