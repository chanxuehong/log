package log

import "io"

// NoopLogger no operation Logger
type NoopLogger struct{}

// Fatal impl Logger Fatal
func (l *NoopLogger) Fatal(msg string, fields ...interface{}) {
}

// Error impl Logger Error
func (l *NoopLogger) Error(msg string, fields ...interface{}) {
}

// Warn impl Logger Warn
func (l *NoopLogger) Warn(msg string, fields ...interface{}) {
}

// Info impl Logger Info
func (l *NoopLogger) Info(msg string, fields ...interface{}) {
}

// Debug impl Logger Debug
func (l *NoopLogger) Debug(msg string, fields ...interface{}) {
}

// Output impl Logger Output
func (l *NoopLogger) Output(calldepth int, level Level, msg string, fields ...interface{}) {
}

// WithField impl Logger WithField
func (l *NoopLogger) WithField(key string, value interface{}) Logger {
	return l
}

// WithFields impl Logger WithFields
func (l *NoopLogger) WithFields(fields ...interface{}) Logger {
	return l
}

// SetFormatter impl Logger SetFormatter
func (l *NoopLogger) SetFormatter(Formatter) {
}

// SetOutput impl Logger SetOutput
func (l *NoopLogger) SetOutput(io.Writer) {
}

// SetLevel impl Logger SetLevel
func (l *NoopLogger) SetLevel(Level) error {
	return nil
}

// SetLevelString impl Logger SetLevelString
func (l *NoopLogger) SetLevelString(string) error {
	return nil
}
