package labels

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ManagedBy      = "app.kubernetes.io/managed-by"
	managedByValue = "glasskube"
)

func IsManaged(obj metav1.Object) bool {
	labels := obj.GetLabels()
	return labels != nil && labels[ManagedBy] == managedByValue
}

func SetManaged(obj metav1.Object) {
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[ManagedBy] = managedByValue
	obj.SetLabels(labels)
}

func IsManagedOrNotExists(c client.Client, ctx context.Context, obj client.Object) (bool, error) {
	if actual, ok := obj.DeepCopyObject().(client.Object); !ok {
		return false, fmt.Errorf("can not deep-copy %v", obj.GetName())
	} else if err := c.Get(ctx, client.ObjectKeyFromObject(obj), actual); err != nil && !errors.IsNotFound(err) {
		return false, err
	} else {
		// err != nil implies that it IsNotFound
		return err != nil || IsManaged(actual), nil
	}
}
