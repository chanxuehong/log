package log

import (
	"sync/atomic"
	"unsafe"
)

type Options []Option

var _defaultOptionsPtr unsafe.Pointer // *Options

func SetDefaultOptions(opts Options) {
	if opts == nil {
		return
	}
	atomic.StorePointer(&_defaultOptionsPtr, unsafe.Pointer(&opts))
}

func getDefaultOptions() Options {
	ptr := (*Options)(atomic.LoadPointer(&_defaultOptionsPtr))
	if ptr != nil {
		return *ptr
	}
	return nil
}
