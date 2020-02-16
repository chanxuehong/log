package log

import "errors"

var (
	_ErrNumberOfFieldsMustNotBeOdd error = errors.New("the number of fields must not be odd")
	_ErrTypeOfFieldKeyMustBeString error = errors.New("the type of field key must be string")
	_ErrFieldKeyMustNotBeEmpty     error = errors.New("the field key must not be empty")
)

func combineFields(m map[string]interface{}, fields []interface{}) (map[string]interface{}, error) {
	if len(fields) == 0 {
		return cloneFields(m), nil
	}
	if len(fields)&1 != 0 {
		return cloneFields(m), _ErrNumberOfFieldsMustNotBeOdd
	}

	m2 := make(map[string]interface{}, 8+len(m)+len(fields)>>1) // 8 is reserved for the standard field
	for k, v := range m {
		m2[k] = v
	}
	for i := 0; i < len(fields); i += 2 {
		k, ok := fields[i].(string)
		if !ok {
			return m2, _ErrTypeOfFieldKeyMustBeString
		}
		if k == "" {
			return m2, _ErrFieldKeyMustNotBeEmpty
		}
		m2[k] = fields[i+1]
	}
	return m2, nil
}

func cloneFields(fields map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}
	m := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		m[k] = v
	}
	return m
}
