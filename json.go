package log

import (
	"bytes"
	"encoding/json"
)

// JSON is a helper function, following is it's code.
//
//  data, _ := json.Marshal(v)
//  return string(data)
func JSON(v interface{}) string {
	buffer := _bufferPool.Get().(*bytes.Buffer)
	defer _bufferPool.Put(buffer)
	buffer.Reset()
	if err := json.NewEncoder(buffer).Encode(v); err != nil {
		return ""
	}
	return string(buffer.Bytes())
}
