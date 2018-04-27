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

func FromContext(ctx context.Context) string {
	if ctx == nil {
		return NewTraceId()
	}
	traceId, ok := ctx.Value(traceIdContextKey{}).(string)
	if ok && traceId != "" {
		return traceId
	}
	return NewTraceId()
}

const TraceIdHeaderKey = "X-Request-Id"

func FromRequest(req *http.Request) string {
	if req == nil {
		return NewTraceId()
	}
	traceId, ok := req.Context().Value(traceIdContextKey{}).(string)
	if ok && traceId != "" {
		return traceId
	}
	traceId = req.Header.Get(TraceIdHeaderKey)
	if traceId != "" {
		return traceId
	}
	return NewTraceId()
}

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
