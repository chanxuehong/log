package log

import (
	"bytes"
	"testing"
)

func TestGetBytesBufferPool(t *testing.T) {
	pool := getBytesBufferPool()

	buffer := pool.Get()
	defer pool.Put(buffer)

	have := buffer.String()
	want := ""
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}

	buffer.WriteString("test")
	pool.Put(buffer)

	buffer = pool.Get()
	have = buffer.String()
	want = "test"
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}
}

func TestSetBytesBufferPool(t *testing.T) {
	pool := getBytesBufferPool()
	if _, ok := pool.(*_BytesBufferPool); !ok {
		t.Error("want type *testBytesBufferPool")
		return
	}

	SetBytesBufferPool(&testBytesBufferPool{})
	defer SetBytesBufferPool(_defaultBytesBufferPool)

	pool = getBytesBufferPool()
	if _, ok := pool.(*testBytesBufferPool); !ok {
		t.Error("want type *testBytesBufferPool")
		return
	}
}

type testBytesBufferPool struct{}

func (*testBytesBufferPool) Get() *bytes.Buffer  { return nil }
func (*testBytesBufferPool) Put(x *bytes.Buffer) {}

func BenchmarkGetBytesBufferPool(b *testing.B) {
	b.ReportAllocs()
	b.SetParallelism(64)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getBytesBufferPool()
		}
	})
}
