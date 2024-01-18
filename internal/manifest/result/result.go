package result

type resultKind int

const (
	ready resultKind = iota
	waiting
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

func (r *ReconcileResult) IsReady() bool {
	return r != nil && r.kind == ready
}

func (r *ReconcileResult) IsWaiting() bool {
	return r != nil && r.kind == waiting
}
