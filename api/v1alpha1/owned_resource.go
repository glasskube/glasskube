package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OwnedResourceRef struct {
	metav1.GroupVersionKind `json:",inline"`
	Name                    string `json:"name"`
	Namespace               string `json:"namespace,omitempty"`
	MarkedForDeletion       bool   `json:"markedForDeletion,omitempty"`
}

func (orr OwnedResourceRef) String() string {
	return orr.GroupVersionKind.String() + " Namespace=" + orr.Namespace + " Name=" + orr.Name
}
