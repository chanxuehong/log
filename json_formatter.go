package log

import (
	"bytes"
	"encoding/json"
)

var JsonFormatter Formatter = &jsonFormatter{}

type jsonFormatter struct{}

func (f *jsonFormatter) Format(entry *Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0, 16<<10))
	}
	var fields map[string]interface{}
	if len(entry.Fields) > 0 {
		prefixFieldClashes(entry.Fields)
		fields = entry.Fields
	} else {
		fields = make(map[string]interface{}, 8)
	}
	fields[fieldKeyTime] = FormatTimeString(entry.Time.In(_beijingLocation))
	fields[fieldKeyLevel] = entry.Level.String()
	fields[fieldKeyTraceId] = entry.TraceId
	fields[fieldKeyLocation] = entry.Location
	fields[fieldKeyMessage] = entry.Message
	if err := json.NewEncoder(buffer).Encode(fields); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
