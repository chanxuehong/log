package trace

import (
	"context"
	"net/http"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := context.Background()

	// empty traceId
	ctx2 := NewContext(ctx, "")
	if ctx != ctx2 {
		t.Error("want equal")
		return
	}

	// serially NewContext with same traceId
	ctx2 = NewContext(ctx, "1234567890")
	ctx3 := NewContext(ctx2, "1234567890")
	if ctx2 != ctx3 {
		t.Error("want equal")
		return
	}

	// parallel NewContext with same traceId
	ctx2 = NewContext(ctx, "1234567890")
	ctx3 = NewContext(ctx, "1234567890")
	if ctx2 == ctx3 {
		t.Error("want not equal")
		return
	}

	// FromContext
	ctx = NewContext(ctx, "1234567890")
	id, ok := FromContext(ctx)
	wantId, wantOk := "1234567890", true
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}
}

func TestFromContext(t *testing.T) {
	// nil context
	var ctx context.Context
	id, ok := FromContext(ctx)
	wantId, wantOk := "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// context does not contain traceId
	ctx = context.Background()
	id, ok = FromContext(ctx)
	wantId, wantOk = "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// context contains empty traceId
	ctx = context.WithValue(context.Background(), traceIdContextKey{}, "")
	id, ok = FromContext(ctx)
	wantId, wantOk = "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// context contains non-empty traceId
	ctx = context.WithValue(context.Background(), traceIdContextKey{}, "1234567890")
	id, ok = FromContext(ctx)
	wantId, wantOk = "1234567890", true
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}
}

func TestFromRequest(t *testing.T) {
	// nil *http.Request
	var req *http.Request
	id, ok := FromRequest(req)
	wantId, wantOk := "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// nil Context() and nil Header
	req = &http.Request{}
	id, ok = FromRequest(req)
	wantId, wantOk = "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// nil Header and non-nil Context()
	{
		// context does not contain traceId
		req2 := req.WithContext(context.Background())
		id, ok = FromRequest(req2)
		wantId, wantOk = "", false
		if id != wantId || ok != wantOk {
			t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
			return
		}

		// context contains empty traceId
		req2 = req.WithContext(context.WithValue(context.Background(), traceIdContextKey{}, ""))
		id, ok = FromRequest(req2)
		wantId, wantOk = "", false
		if id != wantId || ok != wantOk {
			t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
			return
		}

		// context contains non-empty traceId
		req2 = req.WithContext(context.WithValue(context.Background(), traceIdContextKey{}, "1234567890"))
		id, ok = FromRequest(req2)
		wantId, wantOk = "1234567890", true
		if id != wantId || ok != wantOk {
			t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
			return
		}
	}

	// nil Context() and non-nil Header
	{
		// without TraceIdHeaderKey
		header := make(http.Header)
		req.Header = header
		id, ok = FromRequest(req)
		wantId, wantOk = "", false
		if id != wantId || ok != wantOk {
			t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
			return
		}

		// with empty value for TraceIdHeaderKey
		header = make(http.Header)
		header.Set(TraceIdHeaderKey, "")
		req.Header = header
		id, ok = FromRequest(req)
		wantId, wantOk = "", false
		if id != wantId || ok != wantOk {
			t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
			return
		}

		// with valid value for TraceIdHeaderKey
		header = make(http.Header)
		header.Set(TraceIdHeaderKey, "1234567890")
		req.Header = header
		id, ok = FromRequest(req)
		wantId, wantOk = "1234567890", true
		if id != wantId || ok != wantOk {
			t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
			return
		}
	}

	// non-nil Context() and non-nil Header
	{
		header := make(http.Header)
		header.Set(TraceIdHeaderKey, "1234567890-header")
		req.Header = header
		req2 := req.WithContext(context.WithValue(context.Background(), traceIdContextKey{}, "1234567890-context"))
		id, ok = FromRequest(req2)
		wantId, wantOk = "1234567890-context", true
		if id != wantId || ok != wantOk {
			t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
			return
		}
	}
}

func TestFromHeader(t *testing.T) {
	// nil header
	var header http.Header
	id, ok := FromHeader(header)
	wantId, wantOk := "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// without TraceIdHeaderKey
	header = make(http.Header)
	id, ok = FromHeader(header)
	wantId, wantOk = "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// with empty value for TraceIdHeaderKey
	header = make(http.Header)
	header.Set(TraceIdHeaderKey, "")
	id, ok = FromHeader(header)
	wantId, wantOk = "", false
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}

	// with valid value for TraceIdHeaderKey
	header = make(http.Header)
	header.Set(TraceIdHeaderKey, "1234567890")
	id, ok = FromHeader(header)
	wantId, wantOk = "1234567890", true
	if id != wantId || ok != wantOk {
		t.Errorf("have:(%s, %t), want:(%s, %t)", id, ok, wantId, wantOk)
		return
	}
}
