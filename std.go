package log

var _std = _New("")

// Fatal logs a message at FatalLevel on the standard logger.
// For more information see the Logger interface.
func Fatal(msg string, fields ...interface{}) {
	_std.output(1, FatalLevel, msg, fields)
}

// Error logs a message at ErrorLevel on the standard logger.
// For more information see the Logger interface.
func Error(msg string, fields ...interface{}) {
	_std.output(1, ErrorLevel, msg, fields)
}

// Warn logs a message at WarnLevel on the standard logger.
// For more information see the Logger interface.
func Warn(msg string, fields ...interface{}) {
	_std.output(1, WarnLevel, msg, fields)
}

// Info logs a message at InfoLevel on the standard logger.
// For more information see the Logger interface.
func Info(msg string, fields ...interface{}) {
	_std.output(1, InfoLevel, msg, fields)
}

// Debug logs a message at DebugLevel on the standard logger.
// For more information see the Logger interface.
func Debug(msg string, fields ...interface{}) {
	_std.output(1, DebugLevel, msg, fields)
}

// Output logs a message at specified level on the standard logger.
// For more information see the Logger interface.
func Output(calldepth int, level Level, msg string, fields ...interface{}) {
	_std.Output(calldepth+1, level, msg, fields...)
}
