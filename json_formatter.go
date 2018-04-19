package log

import (
	"bytes"
	"encoding/json"
)

type jsonFormatter struct{}

func (f *jsonFormatter) Format(entry *entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0, 4<<10))
	}
	var fields map[string]interface{}
	if len(entry.Fields) > 0 {
		prefixFieldClashes(entry.Fields)
		fields = entry.Fields
	} else {
		fields = make(map[string]interface{}, 8)
	}
	fields["time"] = entry.Time.In(_beijingLocation).Format(TimeFormatLayout)
	fields["level"] = entry.Level.String()
	fields["request_id"] = entry.TraceId
	fields["file_line"] = entry.Location
	fields["msg"] = entry.Message
	if err := json.NewEncoder(buffer).Encode(fields); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
