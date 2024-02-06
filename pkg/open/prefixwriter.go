package open

import "io"

type prefixWriter struct {
	prefix string
	writer io.Writer
}

func (o prefixWriter) Write(data []byte) (int, error) {
	return o.writer.Write(append([]byte(o.prefix), data...))
}

func (o prefixWriter) Close() error {
	if closer, ok := o.writer.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}
