package log

import "encoding/json"

// JSON is a helper function, following is it's code.
//
//  data, _ := json.Marshal(v)
//  return string(data)
func JSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
