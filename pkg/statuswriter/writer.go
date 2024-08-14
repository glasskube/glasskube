package statuswriter

import (
	"fmt"
	"io"
	"os"
)

type writerStatusWriter struct {
	writer    io.Writer
	autoclose bool
}

// SetStatus implements StatusWriter.
func (obj *writerStatusWriter) SetStatus(desc string) {
	// TODO: Handle error returned by fmt.Fprintln
	_, _ = fmt.Fprintln(obj.writer, desc)
}

// Start implements StatusWriter.
func (*writerStatusWriter) Start() {}

// Stop implements StatusWriter.
func (obj *writerStatusWriter) Stop() {
	if obj.autoclose {
		if closer, ok := obj.writer.(io.Closer); ok {
			_ = closer.Close()
		}
	}
}

func Writer(writer io.Writer, autoclose bool) *writerStatusWriter {
	return &writerStatusWriter{writer: writer, autoclose: autoclose}
}

func Stdout() *writerStatusWriter {
	return Writer(os.Stdout, false)
}

func Stderr() *writerStatusWriter {
	return Writer(os.Stderr, false)
}
