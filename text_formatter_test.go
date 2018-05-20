package log

import (
	"reflect"
	"testing"
	"time"
)

func TestTextFormatter_Format(t *testing.T) {
	entry := &Entry{
		Location: "function(file:line)",
		Time:     time.Date(2018, time.May, 20, 8, 20, 30, 666000000, time.UTC),
		Level:    InfoLevel,
		TraceId:  "trace_id_1234567890",
		Message:  "message_1234567890",
		Fields: map[string]interface{}{
			"key1":           "fields_value1",
			"key2":           "fields_value2",
			"key3":           "fields_value3",
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
	want := `time=2018-05-20 16:20:30.666, level=info, request_id=trace_id_1234567890, location=function(file:line), msg=message_1234567890, ` +
		`fields.level=level, fields.location=location, fields.msg=msg, fields.request_id=request_id, fields.time=time, ` +
		`key1=fields_value1, key2=fields_value2, key3=fields_value3` + "\n"
	if string(have) != want {
		t.Errorf("\nhave:%s\nwant:%s", have, want)
		return
	}
}

func TestPrefixFieldClashes(t *testing.T) {
	m := map[string]interface{}{
		"time":           "time",
		"fields.time":    "fields.time",
		"level":          "level",
		"fields.level":   "fields.level",
		"fields.level.2": "fields.level.2",
	}
	prefixFieldClashes(m)
	want := map[string]interface{}{
		"fields.time.2":  "time",
		"fields.time":    "fields.time",
		"fields.level.3": "level",
		"fields.level":   "fields.level",
		"fields.level.2": "fields.level.2",
	}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("\nhave:%v\nwant:%v", m, want)
		return
	}
}
