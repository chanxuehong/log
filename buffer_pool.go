package log

import (
	"bytes"
	"sync"
	"sync/atomic"
	"unsafe"
)

type BytesBufferPool interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

var _BytesBufferPoolPtr unsafe.Pointer = unsafe.Pointer(&_defaultBytesBufferPool)

func getBytesBufferPool() BytesBufferPool {
	ptr := (*BytesBufferPool)(atomic.LoadPointer(&_BytesBufferPoolPtr))
	return *ptr
}

func SetBytesBufferPool(pool BytesBufferPool) {
	if pool == nil {
		return
	}
	atomic.StorePointer(&_BytesBufferPoolPtr, unsafe.Pointer(&pool))
}

var _defaultBytesBufferPool BytesBufferPool = &_BytesBufferPool{
	pool: sync.Pool{
		New: syncPoolNew,
	},
}

func syncPoolNew() interface{} {
	return bytes.NewBuffer(make([]byte, 0, 16<<10))
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
