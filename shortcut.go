package log

import (
	"context"
	"io"
)

// FatalContext is a shortcut for MustFromContext(ctx).Output(1, FatalLevel, msg, fields...)
func FatalContext(ctx context.Context, msg string, fields ...interface{}) {
	MustFromContext(ctx).Output(1, FatalLevel, msg, fields...)
}

// ErrorContext is a shortcut for MustFromContext(ctx).Output(1, ErrorLevel, msg, fields...)
func ErrorContext(ctx context.Context, msg string, fields ...interface{}) {
	MustFromContext(ctx).Output(1, ErrorLevel, msg, fields...)
}

// WarnContext is a shortcut for MustFromContext(ctx).Output(1, WarnLevel, msg, fields...)
func WarnContext(ctx context.Context, msg string, fields ...interface{}) {
	MustFromContext(ctx).Output(1, WarnLevel, msg, fields...)
}

// InfoContext is a shortcut for MustFromContext(ctx).Output(1, InfoLevel, msg, fields...)
func InfoContext(ctx context.Context, msg string, fields ...interface{}) {
	MustFromContext(ctx).Output(1, InfoLevel, msg, fields...)
}

// DebugContext is a shortcut for MustFromContext(ctx).Output(1, DebugLevel, msg, fields...)
func DebugContext(ctx context.Context, msg string, fields ...interface{}) {
	MustFromContext(ctx).Output(1, DebugLevel, msg, fields...)
}

// OutputContext is a shortcut for MustFromContext(ctx).Output(calldepth+1, level, msg, fields...)
func OutputContext(ctx context.Context, calldepth int, level Level, msg string, fields ...interface{}) {
	MustFromContext(ctx).Output(calldepth+1, level, msg, fields...)
}

// WithFieldContext is a shortcut for MustFromContext(ctx).WithField(key, value)
func WithFieldContext(ctx context.Context, key string, value interface{}) Logger {
	return MustFromContext(ctx).WithField(key, value)
}

// WithFieldsContext is a shortcut for MustFromContext(ctx).WithFields(fields...)
func WithFieldsContext(ctx context.Context, fields ...interface{}) Logger {
	return MustFromContext(ctx).WithFields(fields...)
}

// SetFormatterContext is a shortcut for MustFromContext(ctx).SetFormatter(formatter)
func SetFormatterContext(ctx context.Context, formatter Formatter) {
	MustFromContext(ctx).SetFormatter(formatter)
}

// SetOutputContext is a shortcut for MustFromContext(ctx).SetOutput(output)
func SetOutputContext(ctx context.Context, output io.Writer) {
	MustFromContext(ctx).SetOutput(output)
}

// SetLevelContext is a shortcut for MustFromContext(ctx).SetLevel(level)
func SetLevelContext(ctx context.Context, level Level) error {
	return MustFromContext(ctx).SetLevel(level)
}

// SetLevelStringContext is a shortcut for MustFromContext(ctx).SetLevelString(str)
func SetLevelStringContext(ctx context.Context, str string) error {
	return MustFromContext(ctx).SetLevelString(str)
}
