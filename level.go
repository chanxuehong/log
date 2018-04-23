package log

import (
	"fmt"
	"strings"
	"sync/atomic"
)

const (
	InvalidLevel Level = iota // InvalidLevel must equal 0
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

const (
	FatalLevelString = "fatal"
	ErrorLevelString = "error"
	WarnLevelString  = "warning"
	InfoLevelString  = "info"
	DebugLevelString = "debug"
)

func isValidLevel(level Level) bool {
	return level >= FatalLevel && level <= DebugLevel
}
func isLevelEnabled(level, loggerLevel Level) bool {
	return loggerLevel >= level
}

type Level uint

func (level Level) String() string {
	switch level {
	case FatalLevel:
		return FatalLevelString
	case ErrorLevel:
		return ErrorLevelString
	case WarnLevel:
		return WarnLevelString
	case InfoLevel:
		return InfoLevelString
	case DebugLevel:
		return DebugLevelString
	default:
		return fmt.Sprintf("unknown_%d", level)
	}
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) error {
	if !isValidLevel(level) {
		return fmt.Errorf("invalid level: %d", level)
	}
	setLevel(level)
	return nil
}

// SetLevelString sets the standard logger level.
func SetLevelString(str string) error {
	var level Level
	switch strings.ToLower(str) {
	case DebugLevelString:
		level = DebugLevel
	case InfoLevelString:
		level = InfoLevel
	case WarnLevelString:
		level = WarnLevel
	case ErrorLevelString:
		level = ErrorLevel
	case FatalLevelString:
		level = FatalLevel
	default:
		return fmt.Errorf("invalid level string: %q", str)
	}
	setLevel(level)
	return nil
}

var _level = uint64(DebugLevel) // default is debug

func setLevel(level Level) {
	atomic.StoreUint64(&_level, uint64(level))
}

func getLevel() Level {
	return Level(atomic.LoadUint64(&_level))
}
