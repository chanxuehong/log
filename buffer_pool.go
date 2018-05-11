package log

import (
	"bytes"
	"sync"
	"sync/atomic"
	"unsafe"
)

type BytesBufferPool interface {
	Get() *bytes.Buffer
	Put(x *bytes.Buffer)
}

func SetBytesBufferPool(pool BytesBufferPool) {
	if pool == nil {
		return
	}
	atomic.StorePointer(&_BytesBufferPoolPtr, unsafe.Pointer(&pool))
}

func getBytesBufferPool() BytesBufferPool {
	ptr := (*BytesBufferPool)(atomic.LoadPointer(&_BytesBufferPoolPtr))
	return *ptr
}

var _BytesBufferPoolPtr unsafe.Pointer // *BytesBufferPool

func init() {
	SetBytesBufferPool(_NewBytesBufferPool())
}

func _NewBytesBufferPool() BytesBufferPool {
	return &_BytesBufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 16<<10))
			},
		},
	}
}

type _BytesBufferPool struct {
	pool sync.Pool
}

func (p *_BytesBufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

func (p *_BytesBufferPool) Put(x *bytes.Buffer) {
	if x == nil {
		return
	}
	p.pool.Put(x)
}
