package log

import (
	"fmt"
	"strings"
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
