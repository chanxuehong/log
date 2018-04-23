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

func isValidLevel(level Level) bool {
	return level >= FatalLevel && level <= DebugLevel
}
func isLevelEnabled(level, loggerLevel Level) bool {
	return loggerLevel >= level
}

const (
	FatalLevelString = "fatal"
	ErrorLevelString = "error"
	WarnLevelString  = "warning"
	InfoLevelString  = "info"
	DebugLevelString = "debug"
)

func parseLevelString(str string) (level Level, ok bool) {
	switch strings.ToLower(str) {
	case DebugLevelString:
		return DebugLevel, true
	case InfoLevelString:
		return InfoLevel, true
	case WarnLevelString:
		return WarnLevel, true
	case ErrorLevelString:
		return ErrorLevel, true
	case FatalLevelString:
		return FatalLevel, true
	default:
		return InvalidLevel, false
	}
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
	level, ok := parseLevelString(str)
	if !ok {
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
