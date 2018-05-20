package trace

import (
	"context"
	"net/http"

	"github.com/chanxuehong/uuid"
)

type Tracer interface {
	TraceId() string
}

func NewTraceId() string {
	return string(uuid.NewV1().HexEncode())
}

type traceIdContextKey struct{}

func NewContext(ctx context.Context, traceId string) context.Context {
	if traceId == "" {
		return ctx
	}
	if ctx == nil {
		return context.WithValue(context.Background(), traceIdContextKey{}, traceId)
	}
	if v, ok := ctx.Value(traceIdContextKey{}).(string); ok && v == traceId {
		return ctx
	}
	return context.WithValue(ctx, traceIdContextKey{}, traceId)
}

func FromContext(ctx context.Context) (traceId string, ok bool) {
	if ctx == nil {
		return "", false
	}
	traceId, ok = ctx.Value(traceIdContextKey{}).(string)
	if ok && traceId != "" {
		return traceId, true
	}
	return "", false
}

func FromRequest(req *http.Request) (traceId string, ok bool) {
	if req == nil {
		return "", false
	}
	traceId, ok = FromContext(req.Context())
	if ok {
		return traceId, true
	}
	return FromHeader(req.Header)
}

const TraceIdHeaderKey = "X-Request-Id"

func FromHeader(header http.Header) (traceId string, ok bool) {
	traceId = header.Get(TraceIdHeaderKey)
	if traceId != "" {
		return traceId, true
	}
	return "", false
}
