package log

import (
	"context"
	"net/http"
)

type loggerContextKey struct{}

var _loggerContextKey loggerContextKey

func NewContext(ctx context.Context, logger Logger) context.Context {
	if logger == nil {
		return ctx
	}
	if ctx == nil {
		return context.WithValue(context.Background(), _loggerContextKey, logger)
	}
	if value, ok := ctx.Value(_loggerContextKey).(Logger); ok && value == logger {
		return ctx
	}
	return context.WithValue(ctx, _loggerContextKey, logger)
}

func FromContext(ctx context.Context) (lg Logger, ok bool) {
	if ctx == nil {
		return nil, false
	}
	lg, ok = ctx.Value(_loggerContextKey).(Logger)
	return
}

func MustFromContext(ctx context.Context) Logger {
	lg, ok := FromContext(ctx)
	if !ok {
		panic("log: failed to get from context.Context")
	}
	return lg
}

func FromContextOrNew(ctx context.Context, new func() Logger) (lg Logger, isNew bool) {
	lg, ok := FromContext(ctx)
	if ok {
		return lg, false
	}
	if new != nil {
		return new(), true
	}
	return New(), true
}

func FromRequest(req *http.Request) (lg Logger, ok bool) {
	if req == nil {
		return nil, false
	}
	return FromContext(req.Context())
}

func MustFromRequest(req *http.Request) Logger {
	lg, ok := FromRequest(req)
	if !ok {
		panic("log: failed to get from http.Request")
	}
	return lg
}

func FromRequestOrNew(req *http.Request, new func() Logger) (lg Logger, isNew bool) {
	lg, ok := FromRequest(req)
	if ok {
		return lg, false
	}
	if new != nil {
		return new(), true
	}
	return New(), true
}
