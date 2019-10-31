package log

import (
	"bytes"
	"testing"
)

func Test_SetBytesBufferPool_GetBytesBufferPool(t *testing.T) {
	defer SetBytesBufferPool(_defaultBytesBufferPool)

	// get default BytesBufferPool
	pool := getBytesBufferPool()
	if _, ok := pool.(*bytesBufferPool); !ok {
		t.Error("want type *bytesBufferPool")
		return
	}

	// SetBytesBufferPool with nil
	pool = nil
	SetBytesBufferPool(pool)
	pool = getBytesBufferPool()
	if _, ok := pool.(*bytesBufferPool); !ok {
		t.Error("want type *bytesBufferPool")
		return
	}

	// SetBytesBufferPool with non-nil
	SetBytesBufferPool(&testBytesBufferPool{})
	pool = getBytesBufferPool()
	if _, ok := pool.(*testBytesBufferPool); !ok {
		t.Error("want type *testBytesBufferPool")
		return
	}

	// SetBytesBufferPool with nil
	pool = nil
	SetBytesBufferPool(pool)
	pool = getBytesBufferPool()
	if _, ok := pool.(*testBytesBufferPool); !ok {
		t.Error("want type *testBytesBufferPool")
		return
	}
}

type testBytesBufferPool struct{}

func (*testBytesBufferPool) Get() *bytes.Buffer  { return nil }
func (*testBytesBufferPool) Put(x *bytes.Buffer) {}
