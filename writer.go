package log

import (
	"io"
	"os"
	"sync"
)

var (
	concurrentStdout io.Writer = ConcurrentWriter(os.Stdout)
	concurrentStderr io.Writer = ConcurrentWriter(os.Stderr)
)

func ConcurrentWriter(w io.Writer) io.Writer {
	if w == nil {
		return nil
	}
	return &concurrentWriter{
		Writer: w,
	}
}

type concurrentWriter struct {
	sync.Mutex
	io.Writer
}

func (w *concurrentWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	return w.Write(p)
}
