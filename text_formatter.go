package log

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

var TextFormatter Formatter = &textFormatter{}

type textFormatter struct{}

func (f *textFormatter) Format(entry *Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0, 16<<10))
	}
	if len(entry.Fields) > 0 {
		prefixFieldClashes(entry.Fields)
	}
	f.appendKeyValue(buffer, fieldKeyTime, FormatTimeString(entry.Time.In(_beijingLocation)))
	f.appendKeyValue(buffer, fieldKeyLevel, entry.Level.String())
	f.appendKeyValue(buffer, fieldKeyTraceId, entry.TraceId)
	f.appendKeyValue(buffer, fieldKeyLocation, entry.Location)
	f.appendKeyValue(buffer, fieldKeyMessage, entry.Message)

	keys := make([]string, 0, len(entry.Fields))
	for k := range entry.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := entry.Fields[k]
		f.appendKeyValue(buffer, k, v)
	}
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

var _beijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)

const (
	fieldKeyTime     = "time"
	fieldKeyLevel    = "level"
	fieldKeyTraceId  = "request_id"
	fieldKeyLocation = "location"
	fieldKeyMessage  = "msg"
)

func prefixFieldClashes(data map[string]interface{}) {
	if v, ok := data[fieldKeyTime]; ok {
		data["fields."+fieldKeyTime] = v
		delete(data, fieldKeyTime)
	}
	if v, ok := data[fieldKeyLevel]; ok {
		data["fields."+fieldKeyLevel] = v
		delete(data, fieldKeyLevel)
	}
	if v, ok := data[fieldKeyTraceId]; ok {
		data["fields."+fieldKeyTraceId] = v
		delete(data, fieldKeyTraceId)
	}
	if v, ok := data[fieldKeyLocation]; ok {
		data["fields."+fieldKeyLocation] = v
		delete(data, fieldKeyLocation)
	}
	if v, ok := data[fieldKeyMessage]; ok {
		data["fields."+fieldKeyMessage] = v
		delete(data, fieldKeyMessage)
	}
}

func (f *textFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteString(", ")
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *textFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(stringVal)
}
