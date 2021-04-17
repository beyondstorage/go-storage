package iowrap

import (
	"io"
)

// CallbackWriter will create a new CallbackifyWriter.
func CallbackWriter(w io.Writer, fn func([]byte)) *CallbackifyWriter {
	return &CallbackifyWriter{
		w:  w,
		fn: fn,
	}
}

// CallbackifyWriter will execute callback func in Write.
type CallbackifyWriter struct {
	w  io.Writer
	fn func([]byte)
}

// Write will write into underlying Writer.
func (w *CallbackifyWriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.fn(p[:n])
	return n, err
}
