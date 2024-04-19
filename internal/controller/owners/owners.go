package owners

import (
	"context"
	"errors"

	"github.com/glasskube/glasskube/internal/controller/labels"
	"github.com/glasskube/glasskube/internal/controller/owners/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type OwnerManager struct {
	scheme *runtime.Scheme
}

type OwnerOptions int

const (
	BlockOwnerDeletion OwnerOptions = 1 << iota
	Controller
	DefaultOptions OwnerOptions = 0
)

var (
	ErrNoSuchOwner = errors.New("no such owner")
)

func NewOwnerManager(scheme *runtime.Scheme) *OwnerManager {
	return &OwnerManager{scheme: scheme}
}

func (mgr *OwnerManager) GetScheme() *runtime.Scheme {
	return mgr.scheme
}

func (mgr *OwnerManager) HasOwner(owner client.Object, obj metav1.Object) (bool, error) {
	if _, err := mgr.findOwnerReferenceIndex(owner, obj.GetOwnerReferences()); err != nil && !errors.Is(err, ErrNoSuchOwner) {
		return false, err
	} else if errors.Is(err, ErrNoSuchOwner) {
		return false, nil
	} else {
		return true, nil
	}
}

func (mgr *OwnerManager) HasAnyOwnerOfType(owner client.Object, obj metav1.Object) (bool, error) {
	filtered, err := mgr.OwnersOfType(owner, obj)
	return len(filtered) > 0, err
}

func (mgr *OwnerManager) OwnersOfType(owner client.Object, obj metav1.Object) ([]metav1.OwnerReference, error) {
	ownerGVK, err := utils.GetGVK(mgr.scheme, owner)
	if err != nil {
		return nil, err
	}
	ownerGV := schema.GroupVersionKind(ownerGVK).GroupVersion()
	var filtered []metav1.OwnerReference
	for _, ref := range obj.GetOwnerReferences() {
		if ownerGVK.Kind == ref.Kind {
			if refGV, err := schema.ParseGroupVersion(ref.APIVersion); err != nil {
				return filtered, err
			} else if ownerGV == refGV {
				filtered = append(filtered, ref)
			}
		}
	}
	return filtered, nil
}

func (mgr *OwnerManager) CountOwnersOfType(owner client.Object, obj metav1.Object) (int, error) {
	ownerGVK := owner.GetObjectKind().GroupVersionKind()
	ownerGV := ownerGVK.GroupVersion()
	count := 0
	for _, ref := range obj.GetOwnerReferences() {
		if refGV, err := schema.ParseGroupVersion(ref.APIVersion); err != nil {
			return 0, err
		} else if ownerGV == refGV {
			count = count + 1
		}
	}
	return count, nil
}

func (mgr *OwnerManager) SetOwner(
	owner client.Object,
	obj metav1.Object,
	options OwnerOptions,
) error {
	if options&Controller != 0 {
		if err := controllerutil.SetControllerReference(owner, obj, mgr.scheme); err != nil {
			return err
		}
	} else {
		if err := controllerutil.SetOwnerReference(owner, obj, mgr.scheme); err != nil {
			return err
		}
	}
	references := obj.GetOwnerReferences()
	i, err := mgr.findOwnerReferenceIndex(owner, references)
	if err != nil {
		return err
	}
	ref := &references[i]
	blockOwnerDeletion := options&BlockOwnerDeletion != 0
	ref.BlockOwnerDeletion = &blockOwnerDeletion
	obj.SetOwnerReferences(references)

	return nil
}

// SetOwnerIfManagedOrNotExists ensures that the operator only sets owner references on objects when it also manages them.
func (mgr *OwnerManager) SetOwnerIfManagedOrNotExists(c client.Client, ctx context.Context, owner client.Object, obj client.Object) error {
	if shouldSetOwner, err := labels.IsManagedOrNotExists(c, ctx, obj); err != nil {
		return err
	} else if shouldSetOwner {
		labels.SetManaged(obj)
		return mgr.SetOwner(owner, obj, BlockOwnerDeletion)
	} else {
		return nil
	}
}

func (mgr *OwnerManager) RemoveOwner(owner client.Object, obj metav1.Object) error {
	return controllerutil.RemoveOwnerReference(owner, obj, mgr.scheme)
}

func (mgr *OwnerManager) findOwnerReferenceIndex(owner client.Object, references []metav1.OwnerReference) (int, error) {
	ownerName := owner.GetName()
	ownerGVK := owner.GetObjectKind().GroupVersionKind()
	ownerGV := ownerGVK.GroupVersion()
	for i, ref := range references {
		if ref.Name != ownerName || ref.Kind != ownerGVK.Kind {
			continue
		}
		if refGV, err := schema.ParseGroupVersion(ref.APIVersion); err != nil {
			return -1, err
		} else if ownerGV == refGV {
			return i, nil
		}
	}
	return -1, ErrNoSuchOwner
}
