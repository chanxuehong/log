package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

var TextFormatter Formatter = textFormatter{}

type textFormatter struct{}

func (f textFormatter) Format(entry *Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0, 16<<10))
	}
	f.appendKeyValue(buffer, fieldKeyTime, FormatTimeString(entry.Time.In(_beijingLocation)))
	f.appendKeyValue(buffer, fieldKeyLevel, entry.Level.String())
	f.appendKeyValue(buffer, fieldKeyTraceId, entry.TraceId)
	f.appendKeyValue(buffer, fieldKeyLocation, entry.Location)
	f.appendKeyValue(buffer, fieldKeyMessage, entry.Message)
	if fields := entry.Fields; len(fields) > 0 {
		fixFieldsConflictAndHandleErrorFields(fields)

		keys := make([]string, 0, len(fields))
		for k := range fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := fields[k]
			f.appendKeyValue(buffer, k, v)
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

func (f textFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteString(", ")
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f textFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	var stringVal string
	switch v := value.(type) {
	case string:
		stringVal = v
	case json.RawMessage:
		stringVal = string(v)
	default:
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(stringVal)
}
