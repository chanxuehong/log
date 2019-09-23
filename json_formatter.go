package log

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var JsonFormatter Formatter = jsonFormatter{}

type jsonFormatter struct{}

func (jsonFormatter) Format(entry *Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0, 16<<10))
	}
	var fields map[string]interface{}
	if fields = entry.Fields; len(fields) > 0 {
		prefixFieldClashes(fields)
		for k, v := range fields {
			switch vv := v.(type) {
			case error:
				if vv != nil {
					fields[k] = vv.Error()
				}
			case proto.Message:
				m := jsonpb.Marshaler{OrigName: true}
				var buf bytes.Buffer
				if err := m.Marshal(&buf, vv); err != nil {
					return nil, err
				}
				fields[k] = json.RawMessage(buf.Bytes())
			}
		}
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
