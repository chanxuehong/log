package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
)

func testFatalContextLocation() {
	FatalContext(testWithLoggerContext, "msg")
}
func testErrorContextLocation() {
	ErrorContext(testWithLoggerContext, "msg")
}
func testWarnContextLocation() {
	WarnContext(testWithLoggerContext, "msg")
}
func testInfoContextLocation() {
	InfoContext(testWithLoggerContext, "msg")
}
func testDebugContextLocation() {
	DebugContext(testWithLoggerContext, "msg")
}
func testOutputContextLocation() {
	OutputContext(testWithLoggerContext, 0, ErrorLevel, "msg")
}

var testWithLoggerContext = NewContext(context.Background(), New())

func setWithLoggerContextOptionsToDefault() {
	SetFormatterContext(testWithLoggerContext, TextFormatter)
	SetOutputContext(testWithLoggerContext, ConcurrentStdout)
	SetLevelContext(testWithLoggerContext, DebugLevel)
}

//type locationFormat struct{}
//
//func (locationFormat) Format(entry *Entry) ([]byte, error) {
//	return []byte(entry.Location), nil
//}

func TestLocationContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	// fatal
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, locationFormat{})

		testFatalContextLocation()

		location := buf.String()
		switch {
		case location == "log.testFatalContextLocation(github.com/chanxuehong/log/shortcut_test.go:14)":
		case strings.HasPrefix(location, "log.testFatalContextLocation(") && strings.HasSuffix(location, "/log/shortcut_test.go:14)"):
		default:
			t.Errorf("not expected location: %s", location)
			return
		}
	}

	// error
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, locationFormat{})

		testErrorContextLocation()

		location := buf.String()
		switch {
		case location == "log.testErrorContextLocation(github.com/chanxuehong/log/shortcut_test.go:17)":
		case strings.HasPrefix(location, "log.testErrorContextLocation(") && strings.HasSuffix(location, "/log/shortcut_test.go:17)"):
		default:
			t.Errorf("not expected location: %s", location)
			return
		}
	}

	// warning
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, locationFormat{})

		testWarnContextLocation()

		location := buf.String()
		switch {
		case location == "log.testWarnContextLocation(github.com/chanxuehong/log/shortcut_test.go:20)":
		case strings.HasPrefix(location, "log.testWarnContextLocation(") && strings.HasSuffix(location, "/log/shortcut_test.go:20)"):
		default:
			t.Errorf("not expected location: %s", location)
			return
		}
	}

	// info
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, locationFormat{})

		testInfoContextLocation()

		location := buf.String()
		switch {
		case location == "log.testInfoContextLocation(github.com/chanxuehong/log/shortcut_test.go:23)":
		case strings.HasPrefix(location, "log.testInfoContextLocation(") && strings.HasSuffix(location, "/log/shortcut_test.go:23)"):
		default:
			t.Errorf("not expected location: %s", location)
			return
		}
	}

	// debug
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, locationFormat{})

		testDebugContextLocation()

		location := buf.String()
		switch {
		case location == "log.testDebugContextLocation(github.com/chanxuehong/log/shortcut_test.go:26)":
		case strings.HasPrefix(location, "log.testDebugContextLocation(") && strings.HasSuffix(location, "/log/shortcut_test.go:26)"):
		default:
			t.Errorf("not expected location: %s", location)
			return
		}
	}

	// output
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, locationFormat{})

		testOutputContextLocation()

		location := buf.String()
		switch {
		case location == "log.testOutputContextLocation(github.com/chanxuehong/log/shortcut_test.go:29)":
		case strings.HasPrefix(location, "log.testOutputContextLocation(") && strings.HasSuffix(location, "/log/shortcut_test.go:29)"):
		default:
			t.Errorf("not expected location: %s", location)
			return
		}
	}
}

//type testJsonFormatter struct{}
//
//func (testJsonFormatter) Format(entry *Entry) ([]byte, error) {
//	// ignored Entry.Location
//
//	// check entry.Time
//	t := time.Now()
//	d := t.Sub(entry.Time)
//	if d > time.Millisecond || d < -time.Millisecond {
//		return nil, fmt.Errorf("time mismatch, have:%v, want:%v", entry.Time, t)
//	}
//
//	// ignored entry.Time
//	m := make(map[string]interface{})
//	prefixFieldClashes(entry.Fields)
//	m[fieldKeyTraceId] = entry.TraceId
//	m[fieldKeyLevel] = entry.Level.String()
//	m[fieldKeyMessage] = entry.Message
//	for k, v := range entry.Fields {
//		m[k] = v
//	}
//	return json.Marshal(m)
//}

func TestFatalContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	var buf bytes.Buffer
	SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
	SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

	FatalContext(testWithLoggerContext, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
	data := buf.Bytes()

	var have map[string]interface{}
	if err := json.Unmarshal(data, &have); err != nil {
		t.Error(err.Error())
		return
	}
	want := map[string]interface{}{
		fieldKeyTraceId: "",
		fieldKeyLevel:   FatalLevelString,
		fieldKeyMessage: "fatal-msg",
		"field1-key":    "field1-value",
		"field2-key":    "field2-value",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave:%v\nwant:%v", have, want)
		return
	}
}

func TestErrorContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	var buf bytes.Buffer
	SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
	SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

	ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
	data := buf.Bytes()

	var have map[string]interface{}
	if err := json.Unmarshal(data, &have); err != nil {
		t.Error(err.Error())
		return
	}
	want := map[string]interface{}{
		fieldKeyTraceId: "",
		fieldKeyLevel:   ErrorLevelString,
		fieldKeyMessage: "error-msg",
		"field1-key":    "field1-value",
		"field2-key":    "field2-value",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave:%v\nwant:%v", have, want)
		return
	}
}

func TestWarnContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	var buf bytes.Buffer
	SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
	SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

	WarnContext(testWithLoggerContext, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
	data := buf.Bytes()

	var have map[string]interface{}
	if err := json.Unmarshal(data, &have); err != nil {
		t.Error(err.Error())
		return
	}
	want := map[string]interface{}{
		fieldKeyTraceId: "",
		fieldKeyLevel:   WarnLevelString,
		fieldKeyMessage: "warning-msg",
		"field1-key":    "field1-value",
		"field2-key":    "field2-value",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave:%v\nwant:%v", have, want)
		return
	}
}

func TestInfoContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	var buf bytes.Buffer
	SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
	SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

	InfoContext(testWithLoggerContext, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
	data := buf.Bytes()

	var have map[string]interface{}
	if err := json.Unmarshal(data, &have); err != nil {
		t.Error(err.Error())
		return
	}
	want := map[string]interface{}{
		fieldKeyTraceId: "",
		fieldKeyLevel:   InfoLevelString,
		fieldKeyMessage: "info-msg",
		"field1-key":    "field1-value",
		"field2-key":    "field2-value",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave:%v\nwant:%v", have, want)
		return
	}
}

func TestDebugContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	var buf bytes.Buffer
	SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
	SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

	DebugContext(testWithLoggerContext, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
	data := buf.Bytes()

	var have map[string]interface{}
	if err := json.Unmarshal(data, &have); err != nil {
		t.Error(err.Error())
		return
	}
	want := map[string]interface{}{
		fieldKeyTraceId: "",
		fieldKeyLevel:   DebugLevelString,
		fieldKeyMessage: "debug-msg",
		"field1-key":    "field1-value",
		"field2-key":    "field2-value",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave:%v\nwant:%v", have, want)
		return
	}
}

func TestOutputContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	// fatal
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

		OutputContext(testWithLoggerContext, 0, FatalLevel, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   FatalLevelString,
			fieldKeyMessage: "fatal-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}
	// error
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

		OutputContext(testWithLoggerContext, 0, ErrorLevel, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   ErrorLevelString,
			fieldKeyMessage: "error-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}
	// warning
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

		OutputContext(testWithLoggerContext, 0, WarnLevel, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   WarnLevelString,
			fieldKeyMessage: "warning-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}
	// info
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

		OutputContext(testWithLoggerContext, 0, InfoLevel, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   InfoLevelString,
			fieldKeyMessage: "info-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}
	// debug
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

		OutputContext(testWithLoggerContext, 0, DebugLevel, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   DebugLevelString,
			fieldKeyMessage: "debug-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}
	// fatal - 1
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

		OutputContext(testWithLoggerContext, 0, FatalLevel-1, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()
		if len(data) != 0 {
			t.Errorf("want empty, but now is: %s", data)
			return
		}
	}
	// debug + 1
	{
		var buf bytes.Buffer
		SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
		SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

		OutputContext(testWithLoggerContext, 0, DebugLevel+1, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()
		if len(data) != 0 {
			t.Errorf("want empty, but now is: %s", data)
			return
		}
	}
}

func TestWithFieldContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	var buf bytes.Buffer
	SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
	SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

	{
		l := WithFieldContext(testWithLoggerContext, "field100-key", "field100-value")
		l.Error("error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   ErrorLevelString,
			fieldKeyMessage: "error-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
			"field100-key":  "field100-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}

	buf.Reset()

	// WithField cannot affect original logger
	{
		ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   ErrorLevelString,
			fieldKeyMessage: "error-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}
}

func TestWithFieldsContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	var buf bytes.Buffer
	SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
	SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

	{
		l := WithFieldsContext(testWithLoggerContext, "field100-key", "field100-value", "field200-key", "field200-value")
		l.Error("error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   ErrorLevelString,
			fieldKeyMessage: "error-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
			"field100-key":  "field100-value",
			"field200-key":  "field200-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}

	buf.Reset()

	// WithFields cannot affect original logger
	{
		ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
		data := buf.Bytes()

		var have map[string]interface{}
		if err := json.Unmarshal(data, &have); err != nil {
			t.Error(err.Error())
			return
		}
		want := map[string]interface{}{
			fieldKeyTraceId: "",
			fieldKeyLevel:   ErrorLevelString,
			fieldKeyMessage: "error-msg",
			"field1-key":    "field1-value",
			"field2-key":    "field2-value",
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave:%v\nwant:%v", have, want)
			return
		}
	}
}

func TestSetFormatterContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	// default
	{
		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().formatter
		want := TextFormatter
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	// nil Formatter
	{
		var formatter Formatter
		SetFormatterContext(testWithLoggerContext, formatter)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().formatter
		want := TextFormatter
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	// non-nil Formatter
	{
		var formatter = testJsonFormatter{}
		SetFormatterContext(testWithLoggerContext, formatter)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().formatter
		want := formatter
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
}

func TestSetOutputContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	// default
	{
		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().output
		want := ConcurrentStdout
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	// nil output
	{
		var output io.Writer
		SetOutputContext(testWithLoggerContext, output)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().output
		want := ConcurrentStdout
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	// non-nil output
	{
		var buf bytes.Buffer
		var output io.Writer = ConcurrentWriter(&buf)
		SetOutputContext(testWithLoggerContext, output)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().output
		want := output
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
}

func TestSetLevelContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	{
		SetLevelContext(testWithLoggerContext, FatalLevel)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := FatalLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, ErrorLevel)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := ErrorLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, WarnLevel)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := WarnLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, InfoLevel)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, DebugLevel)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := DebugLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, InfoLevel)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, FatalLevel-1)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, DebugLevel+1)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelContext(testWithLoggerContext, InfoLevel)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
}

func TestSetLevelStringContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	{
		SetLevelStringContext(testWithLoggerContext, FatalLevelString)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := FatalLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, ErrorLevelString)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := ErrorLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, WarnLevelString)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := WarnLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, InfoLevelString)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, DebugLevelString)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := DebugLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, InfoLevelString)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, "panic")

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, "trace")

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
	{
		SetLevelStringContext(testWithLoggerContext, InfoLevelString)

		have := MustFromContext(testWithLoggerContext).(*logger).getOptions().level
		want := InfoLevel
		if have != want {
			t.Errorf("have:%v, want:%v", have, want)
			return
		}
	}
}

func TestLeveledOutputContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	// fatal
	{
		SetLevelContext(testWithLoggerContext, FatalLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, FatalLevel, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, ErrorLevel, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, WarnLevel, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, InfoLevel, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, DebugLevel, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// error
	{
		SetLevelContext(testWithLoggerContext, ErrorLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, FatalLevel, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, ErrorLevel, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, WarnLevel, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, InfoLevel, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, DebugLevel, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// warning
	{
		SetLevelContext(testWithLoggerContext, WarnLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, FatalLevel, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, ErrorLevel, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, WarnLevel, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   WarnLevelString,
				fieldKeyMessage: "warning-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, InfoLevel, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, DebugLevel, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// info
	{
		SetLevelContext(testWithLoggerContext, InfoLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, FatalLevel, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, ErrorLevel, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, WarnLevel, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   WarnLevelString,
				fieldKeyMessage: "warning-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, InfoLevel, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   InfoLevelString,
				fieldKeyMessage: "info-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, DebugLevel, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// debug
	{
		SetLevelContext(testWithLoggerContext, DebugLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, FatalLevel, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, ErrorLevel, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, WarnLevel, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   WarnLevelString,
				fieldKeyMessage: "warning-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, InfoLevel, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   InfoLevelString,
				fieldKeyMessage: "info-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			OutputContext(testWithLoggerContext, 0, DebugLevel, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   DebugLevelString,
				fieldKeyMessage: "debug-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
	}
}

func TestLeveledPrintContext(t *testing.T) {
	defer setWithLoggerContextOptionsToDefault()

	// fatal
	{
		SetLevelContext(testWithLoggerContext, FatalLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			FatalContext(testWithLoggerContext, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			WarnContext(testWithLoggerContext, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			InfoContext(testWithLoggerContext, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			DebugContext(testWithLoggerContext, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// error
	{
		SetLevelContext(testWithLoggerContext, ErrorLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			FatalContext(testWithLoggerContext, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			WarnContext(testWithLoggerContext, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			InfoContext(testWithLoggerContext, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			DebugContext(testWithLoggerContext, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// warning
	{
		SetLevelContext(testWithLoggerContext, WarnLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			FatalContext(testWithLoggerContext, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			WarnContext(testWithLoggerContext, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   WarnLevelString,
				fieldKeyMessage: "warning-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			InfoContext(testWithLoggerContext, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			DebugContext(testWithLoggerContext, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// info
	{
		SetLevelContext(testWithLoggerContext, InfoLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			FatalContext(testWithLoggerContext, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			WarnContext(testWithLoggerContext, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   WarnLevelString,
				fieldKeyMessage: "warning-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			InfoContext(testWithLoggerContext, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   InfoLevelString,
				fieldKeyMessage: "info-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			DebugContext(testWithLoggerContext, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()
			if len(data) != 0 {
				t.Errorf("want empty, but now is: %s", data)
				return
			}
		}
	}

	// debug
	{
		SetLevelContext(testWithLoggerContext, DebugLevel)

		// fatal
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			FatalContext(testWithLoggerContext, "fatal-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   FatalLevelString,
				fieldKeyMessage: "fatal-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// error
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			ErrorContext(testWithLoggerContext, "error-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   ErrorLevelString,
				fieldKeyMessage: "error-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// warning
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			WarnContext(testWithLoggerContext, "warning-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   WarnLevelString,
				fieldKeyMessage: "warning-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// info
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			InfoContext(testWithLoggerContext, "info-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   InfoLevelString,
				fieldKeyMessage: "info-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
		// debug
		{
			var buf bytes.Buffer
			SetOutputContext(testWithLoggerContext, ConcurrentWriter(&buf))
			SetFormatterContext(testWithLoggerContext, testJsonFormatter{})

			DebugContext(testWithLoggerContext, "debug-msg", "field1-key", "field1-value", "field2-key", "field2-value")
			data := buf.Bytes()

			var have map[string]interface{}
			if err := json.Unmarshal(data, &have); err != nil {
				t.Error(err.Error())
				return
			}
			want := map[string]interface{}{
				fieldKeyTraceId: "",
				fieldKeyLevel:   DebugLevelString,
				fieldKeyMessage: "debug-msg",
				"field1-key":    "field1-value",
				"field2-key":    "field2-value",
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("\nhave:%v\nwant:%v", have, want)
				return
			}
		}
	}
}
