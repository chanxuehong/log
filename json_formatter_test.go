package log

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestJsonFormatter_Format(t *testing.T) {
	entry := &Entry{
		Location: "function(file:line)",
		Time:     time.Date(2018, time.May, 20, 8, 20, 30, 666000000, time.UTC),
		Level:    InfoLevel,
		TraceId:  "trace_id_123456789",
		Message:  "message_123456789",
		Fields: map[string]interface{}{
			"key1":           "fields_value1",
			"key2":           "fields_value2",
			"key3":           testError{},         // error
			"key4":           testContextError1{}, // error with ErrorContext
			"key5":           testContextError2{}, // error with ErrorContextJSON
			"key6":           testContextError3{}, // error with ErrorContext and ErrorContextJSON
			fieldKeyTime:     "time",
			fieldKeyLevel:    "level",
			fieldKeyTraceId:  "request_id",
			fieldKeyLocation: "location",
			fieldKeyMessage:  "msg",
		},
		Buffer: nil,
	}
	data, err := JsonFormatter.Format(entry)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Error("want end with '\n'")
		return
	}
	var have map[string]interface{}
	if err = json.Unmarshal(data, &have); err != nil {
		t.Error(err.Error())
		return
	}
	want := map[string]interface{}{
		"time":             "2018-05-20 16:20:30.666",
		"level":            "info",
		"request_id":       "trace_id_123456789",
		"location":         "function(file:line)",
		"msg":              "message_123456789",
		"field.level":      "level",
		"field.location":   "location",
		"field.msg":        "msg",
		"field.request_id": "request_id",
		"field.time":       "time",
		"key1":             "fields_value1",
		"key2":             "fields_value2",
		"key3":             "test_error_123456789", // error
		"key4":             "context_error1_error_123456789",
		"key4_context":     "context_error1_context_123456789",
		"key5":             "context_error2_error_123456789",
		"key5_context":     map[string]interface{}{"key": "context_error2_context_json_123456789"},
		"key6":             "context_error3_error_123456789",
		"key6_context":     map[string]interface{}{"key": "context_error3_context_json_123456789"},
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave:%v\nwant:%v", have, want)
		return
	}
}
