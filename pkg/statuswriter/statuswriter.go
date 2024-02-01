package statuswriter

type StatusWriter interface {
	SetStatus(desc string)
	Start()
	Stop()
}
