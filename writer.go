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
		w: w,
	}
}

type concurrentWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (w *concurrentWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.Write(p)
}
