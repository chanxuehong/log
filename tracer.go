package log

import "github.com/chanxuehong/log/trace"

var _ trace.Tracer = (*logger)(nil)

func (l *logger) TraceId() string {
	return l.getOptions().traceId
}
