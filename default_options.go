package log

import (
	"sync/atomic"
	"unsafe"
)

type OptionsFunc func() []Option

var _defaultOptionsFuncPtr unsafe.Pointer // *OptionsFunc

func SetDefaultOptionsFunc(fn OptionsFunc) {
	if fn == nil {
		return
	}
	atomic.StorePointer(&_defaultOptionsFuncPtr, unsafe.Pointer(&fn))
}

func getDefaultOptionsFunc() OptionsFunc {
	ptr := (*OptionsFunc)(atomic.LoadPointer(&_defaultOptionsFuncPtr))
	if ptr != nil {
		return *ptr
	}
	return nil
}
