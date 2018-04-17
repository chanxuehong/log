package log

import (
	"strconv"
	"sync/atomic"
)

const (
	ErrorLevel Level = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

type Level uint

func (level Level) String() string {
	switch level {
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warning"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	default:
		return "unknown_" + strconv.FormatUint(uint64(level), 10)
	}
}

var _level = uint64(DebugLevel)

func SetLevel(level Level) {
	atomic.StoreUint64(&_level, uint64(level))
}

func getLevel() Level {
	return Level(atomic.LoadUint64(&_level))
}
