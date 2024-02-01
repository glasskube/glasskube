package statuswriter

type noopStatusWriter struct{}

// SetStatus implements StatusWriter.
func (*noopStatusWriter) SetStatus(desc string) {}

// Start implements StatusWriter.
func (*noopStatusWriter) Start() {}

// Stop implements StatusWriter.
func (*noopStatusWriter) Stop() {}

func Noop() *noopStatusWriter {
	return &noopStatusWriter{}
}
