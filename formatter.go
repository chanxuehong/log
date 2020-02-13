package log

import (
	"encoding/json"
	"strconv"
	"time"
)

var _beijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)

const (
	fieldKeyTime     = "time"
	fieldKeyLevel    = "level"
	fieldKeyTraceId  = "request_id"
	fieldKeyLocation = "location"
	fieldKeyMessage  = "msg"
)

var stdFieldKeys = []string{
	fieldKeyTime,
	fieldKeyLevel,
	fieldKeyTraceId,
	fieldKeyLocation,
	fieldKeyMessage,
}

func fixFieldsConflict(fields map[string]interface{}, fieldKeys []string) {
	for _, fieldKey := range stdFieldKeys {
		fieldValue, ok := fields[fieldKey]
		if !ok {
			continue
		}
		delete(fields, fieldKey)
		newKey := "field." + fieldKey
		for key, i := newKey, 2; ; i++ {
			if _, ok = fields[key]; !ok {
				fields[key] = fieldValue
				break
			}
			key = newKey + "." + strconv.Itoa(i)
		}
	}
	for _, fieldKey := range fieldKeys {
		fieldValue, ok := fields[fieldKey]
		if !ok {
			continue
		}
		delete(fields, fieldKey)
		newKey := "field." + fieldKey
		for key, i := newKey, 2; ; i++ {
			if _, ok = fields[key]; !ok {
				fields[key] = fieldValue
				break
			}
			key = newKey + "." + strconv.Itoa(i)
		}
	}
}

func fixFieldsConflictAndHandleErrorFields(fields map[string]interface{}) {
	var (
		errorContextFields map[string]interface{}
		errorContextKeys   []string
	)
	for k, v := range fields {
		errorValue, ok := v.(error)
		if !ok {
			continue
		}
		fields[k] = errorValue.Error()
		switch vv := v.(type) {
		case interface{ ErrorContextJSON() json.RawMessage }:
			if errorContextFields == nil {
				errorContextFields = make(map[string]interface{}, 8)
			}
			errorContextKey := k + "_context"
			errorContextFields[errorContextKey] = vv.ErrorContextJSON()
			errorContextKeys = append(errorContextKeys, errorContextKey)
		case interface{ ErrorContext() string }:
			if errorContextFields == nil {
				errorContextFields = make(map[string]interface{}, 8)
			}
			errorContextKey := k + "_context"
			errorContextFields[errorContextKey] = vv.ErrorContext()
			errorContextKeys = append(errorContextKeys, errorContextKey)
		}
	}
	fixFieldsConflict(fields, errorContextKeys)
	for k, v := range errorContextFields {
		fields[k] = v
	}
}
