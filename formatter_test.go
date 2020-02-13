package log

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFixFieldsConflict(t *testing.T) {
	m := map[string]interface{}{
		"request_id":    "request_id",
		"time":          "time",
		"field.time":    "field.time",
		"level":         "level",
		"field.level":   "field.level",
		"field.level.2": "field.level.2",
	}
	fixFieldsConflict(m, []string{"request_id", "field.level", "field.level.3"})
	want := map[string]interface{}{
		"field.request_id":    "request_id",
		"field.time.2":        "time",
		"field.time":          "field.time",
		"field.field.level.3": "level",
		"field.field.level":   "field.level",
		"field.level.2":       "field.level.2",
	}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("\nhave:%v\nwant:%v", m, want)
		return
	}
}

func TestFixFieldsConflictAndHandleErrorFields(t *testing.T) {
	var nilErr error
	fields := map[string]interface{}{
		"user_id":                123456,                             // integer
		"name":                   "jack",                             // string
		"nil_error":              nilErr,                             // nil error
		"error":                  testError{},                        // error
		"context_error1":         testContextError1{},                // error with ErrorContext
		"context_error2":         testContextError2{},                // error with ErrorContextJSON
		"context_error3":         testContextError3{},                // error with ErrorContext and ErrorContextJSON
		"context_not_error":      testContextWithoutError{X: "test"}, // not error with ErrorContext and ErrorContextJSON
		"context_error1_context": "context_error1_context_value",     // conflict with context_error1.context
	}
	fixFieldsConflictAndHandleErrorFields(fields)
	want := map[string]interface{}{
		"user_id":                      123456,
		"name":                         "jack",
		"nil_error":                    nil,
		"error":                        "test_error_123456789",
		"context_error1":               "context_error1_error_123456789",
		"context_error1_context":       "context_error1_context_123456789",
		"context_error2":               "context_error2_error_123456789",
		"context_error2_context":       json.RawMessage(`{"key":"context_error2_context_json_123456789"}`),
		"context_error3":               "context_error3_error_123456789",
		"context_error3_context":       json.RawMessage(`{"key":"context_error3_context_json_123456789"}`),
		"context_not_error":            testContextWithoutError{X: "test"},
		"field.context_error1_context": "context_error1_context_value",
	}
	if !reflect.DeepEqual(fields, want) {
		t.Errorf("\nhave:%v\nwant:%v", fields, want)
		return
	}
}

type testError struct{}

func (testError) Error() string { return "test_error_123456789" }

type testContextError1 struct{}

func (testContextError1) Error() string { return "context_error1_error_123456789" }

func (testContextError1) ErrorContext() string { return "context_error1_context_123456789" }

type testContextError2 struct{}

func (testContextError2) Error() string { return "context_error2_error_123456789" }

func (testContextError2) ErrorContextJSON() json.RawMessage {
	return []byte(`{"key":"context_error2_context_json_123456789"}`)
}

type testContextError3 struct{}

func (testContextError3) Error() string { return "context_error3_error_123456789" }

func (testContextError3) ErrorContext() string { return "context_error3_context_123456789" }

func (testContextError3) ErrorContextJSON() json.RawMessage {
	return []byte(`{"key":"context_error3_context_json_123456789"}`)
}

type testContextWithoutError struct {
	X string `json:"x"`
}

func (testContextWithoutError) ErrorContext() string { return "context_not_error_context_123456789" }

func (testContextWithoutError) ErrorContextJSON() json.RawMessage {
	return []byte(`{"key":"context_not_error_context_json_123456789"}`)
}
