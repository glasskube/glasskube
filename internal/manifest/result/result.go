package result

import "github.com/glasskube/glasskube/api/v1alpha1"

type resultKind int

const (
	ready resultKind = iota
	waiting
	failed
)

type ReconcileResult struct {
	kind           resultKind
	Message        string
	OwnedResources []v1alpha1.OwnedResourceRef
}

func Ready(message string, ownedResources []v1alpha1.OwnedResourceRef) *ReconcileResult {
	return &ReconcileResult{kind: ready, Message: message, OwnedResources: ownedResources}
}

func Waiting(message string, ownedResources []v1alpha1.OwnedResourceRef) *ReconcileResult {
	return &ReconcileResult{kind: waiting, Message: message, OwnedResources: ownedResources}
}

func Failed(message string, ownedResources []v1alpha1.OwnedResourceRef) *ReconcileResult {
	return &ReconcileResult{kind: failed, Message: message, OwnedResources: ownedResources}
}

func (r *ReconcileResult) IsReady() bool {
	return r != nil && r.kind == ready
}

func (r *ReconcileResult) IsWaiting() bool {
	return r != nil && r.kind == waiting
}

func (r *ReconcileResult) IsFailed() bool {
	return r != nil && r.kind == failed
}
