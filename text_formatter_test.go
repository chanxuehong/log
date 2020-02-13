package log

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTextFormatter_Format(t *testing.T) {
	entry := &Entry{
		Location: "function(file:line)",
		Time:     time.Date(2018, time.May, 20, 8, 20, 30, 666777888, time.UTC),
		Level:    InfoLevel,
		TraceId:  "trace_id_123456789",
		Message:  "message_123456789",
		Fields: map[string]interface{}{
			"key1":           "fields_value1",
			"key2":           "fields_value2",
			"key3":           testError{}, // error
			"key4":           json.RawMessage([]byte(`{"code":0,"msg":""}`)),
			"key5":           testContextError1{}, // error with ErrorContext
			"key6":           testContextError2{}, // error with ErrorContextJSON
			"key7":           testContextError3{}, // error with ErrorContext and ErrorContextJSON
			fieldKeyTime:     "time",
			fieldKeyLevel:    "level",
			fieldKeyTraceId:  "request_id",
			fieldKeyLocation: "location",
			fieldKeyMessage:  "msg",
		},
		Buffer: nil,
	}
	have, err := TextFormatter.Format(entry)
	if err != nil {
		t.Error(err.Error())
		return
	}
	want := `time=2018-05-20 16:20:30.666, level=info, request_id=trace_id_123456789, location=function(file:line), msg=message_123456789, ` +
		`field.level=level, field.location=location, field.msg=msg, field.request_id=request_id, field.time=time, ` +
		`key1=fields_value1, key2=fields_value2, key3=test_error_123456789, key4={"code":0,"msg":""}, ` +
		`key5=context_error1_error_123456789, key5_context=context_error1_context_123456789, ` +
		`key6=context_error2_error_123456789, key6_context={"key":"context_error2_context_json_123456789"}, ` +
		`key7=context_error3_error_123456789, key7_context={"key":"context_error3_context_json_123456789"}` + "\n"
	if string(have) != want {
		t.Errorf("\nhave:%s\nwant:%s", have, want)
		return
	}
}
