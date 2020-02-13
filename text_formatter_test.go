package log

import (
	"encoding/json"
	"reflect"
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
			"key3":           &testError{X: "123456789"}, // error
			"key4":           json.RawMessage([]byte(`{"code":0,"msg":""}`)),
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
		`key1=fields_value1, key2=fields_value2, key3=test_error_123456789, key4={"code":0,"msg":""}` + "\n"
	if string(have) != want {
		t.Errorf("\nhave:%s\nwant:%s", have, want)
		return
	}
}

func TestPrefixFieldClashes(t *testing.T) {
	m := map[string]interface{}{
		"time":          "time",
		"field.time":    "field.time",
		"level":         "level",
		"field.level":   "field.level",
		"field.level.2": "field.level.2",
	}
	prefixFieldClashes(m)
	want := map[string]interface{}{
		"field.time.2":  "time",
		"field.time":    "field.time",
		"field.level.3": "level",
		"field.level":   "field.level",
		"field.level.2": "field.level.2",
	}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("\nhave:%v\nwant:%v", m, want)
		return
	}
}
