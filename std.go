package log

import "os"

var _std = New("")

func Fatal(msg string, fields ...interface{}) {
	_std.Output(1, FatalLevel, msg, fields...)
	os.Exit(1)
}

func Error(msg string, fields ...interface{}) {
	_std.Output(1, ErrorLevel, msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	_std.Output(1, WarnLevel, msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	_std.Output(1, InfoLevel, msg, fields...)
}

func Debug(msg string, fields ...interface{}) {
	_std.Output(1, DebugLevel, msg, fields...)
}

func Output(calldepth int, level Level, msg string, fields ...interface{}) {
	_std.Output(calldepth+1, level, msg, fields...)
}
