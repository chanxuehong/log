package log

import (
	"bytes"
	"testing"
)

func TestConcurrentWriter_Write(t *testing.T) {
	w := bytes.NewBuffer(make([]byte, 0, 64))
	cw := ConcurrentWriter(w)
	cw.Write([]byte("1234567890"))
	have := w.String()
	want := "1234567890"
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}
}
