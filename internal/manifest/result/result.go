package result

type resultKind int

const (
	ready resultKind = iota
	waiting
	failed
)

type ReconcileResult struct {
	kind    resultKind
	Message string
}

func Ready(message string) *ReconcileResult {
	return &ReconcileResult{kind: ready, Message: message}
}

func Waiting(message string) *ReconcileResult {
	return &ReconcileResult{kind: waiting, Message: message}
}

func Failed(message string) *ReconcileResult {
	return &ReconcileResult{kind: failed, Message: message}
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
