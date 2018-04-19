package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	TimeFormatLayout = "2006-01-02 15:04:05.000"
	TraceIdHeaderKey = "X-Request-Id"
)

func FromRequest(req *http.Request) Logger {
	v, ok := req.Context().Value(loggerKey{}).(Logger)
	if ok && v != nil {
		return v
	}
	traceId := req.Header.Get(TraceIdHeaderKey)
	if traceId == "" {
		traceId = NewTraceId()
	}
	return New(traceId)
}

type loggerKey struct{}

func FromContext(ctx context.Context) Logger {
	v, ok := ctx.Value(loggerKey{}).(Logger)
	if ok && v != nil {
		return v
	}
	return New(NewTraceId())
}

func NewContext(ctx context.Context, logger Logger) context.Context {
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, loggerKey{}, logger)
}

type Logger interface {
	Fatal(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Output(calldepth int, level Level, msg string, fields ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields ...interface{}) Logger
}

func New(traceId string) Logger { return _New(traceId) }

func _New(traceId string) *logger {
	return &logger{
		traceId:   traceId,
		fields:    nil,
		out:       os.Stdout,
		formatter: &textFormatter{},
	}
}

type logger struct {
	traceId   string
	fields    map[string]interface{}
	out       io.Writer
	formatter formatter
}

type formatter interface {
	Format(entry *entry) ([]byte, error)
}

type entry struct {
	Location string // file:line
	Time     time.Time
	Level    Level
	TraceId  string
	Message  string
	Fields   map[string]interface{}
	Buffer   *bytes.Buffer
}

func (l *logger) Fatal(msg string, fields ...interface{}) {
	l.output(1, FatalLevel, msg, fields)
	os.Exit(1)
}
func (l *logger) Error(msg string, fields ...interface{}) {
	l.output(1, ErrorLevel, msg, fields)
}
func (l *logger) Warn(msg string, fields ...interface{}) {
	l.output(1, WarnLevel, msg, fields)
}
func (l *logger) Info(msg string, fields ...interface{}) {
	l.output(1, InfoLevel, msg, fields)
}
func (l *logger) Debug(msg string, fields ...interface{}) {
	l.output(1, DebugLevel, msg, fields)
}

var _bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4<<10))
	},
}

func (l *logger) Output(calldepth int, level Level, msg string, fields ...interface{}) {
	if !isValidLevel(level) {
		return
	}
	l.output(calldepth+1, level, msg, fields)
	if level == FatalLevel {
		os.Exit(1)
	}
}

func (l *logger) output(calldepth int, level Level, msg string, fields []interface{}) {
	if !isLevelEnabled(level) {
		return
	}

	var location string
	if _, file, line, ok := runtime.Caller(calldepth + 1); ok {
		location = trimFileName(file) + ":" + strconv.Itoa(line)
	} else {
		location = "???"
	}

	var m map[string]interface{}
	if len(fields) == 0 {
		m = l.fields
	} else {
		m2, err := parseFields(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "log: failed to parseFields, error=%v, location=%s\n", err, location)
			return
		}
		m = make(map[string]interface{}, len(l.fields)+len(m2))
		for k, v := range l.fields {
			m[k] = v
		}
		for k, v := range m2 {
			m[k] = v
		}
	}

	buffer := _bufferPool.Get().(*bytes.Buffer)
	defer _bufferPool.Put(buffer)
	buffer.Reset()

	data, err := l.formatter.Format(&entry{
		Location: location,
		Time:     time.Now(),
		Level:    level,
		TraceId:  l.traceId,
		Message:  msg,
		Fields:   m,
		Buffer:   buffer,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to format entry, error=%v, location=%s\n", err, location)
		return
	}
	if _, err = l.out.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to write to log, error=%v, location=%s\n", err, location)
		return
	}
}

func trimFileName(name string) string {
	i := strings.Index(name, "/src/")
	if i < 0 {
		return name
	}
	name = name[i+len("/src/"):]
	i = strings.Index(name, "/vendor/")
	if i < 0 {
		return name
	}
	return name[i+len("/vendor/"):]
}

func (l *logger) WithField(key string, value interface{}) Logger {
	m := make(map[string]interface{}, len(l.fields)+1)
	for k, v := range l.fields {
		m[k] = v
	}
	m[key] = value
	return &logger{
		traceId:   l.traceId,
		fields:    m,
		out:       l.out,
		formatter: l.formatter,
	}
}
func (l *logger) WithFields(fields ...interface{}) Logger {
	if len(fields) == 0 {
		return l
	}
	m, err := parseFields(fields)
	if err != nil {
		var location string
		if _, file, line, ok := runtime.Caller(1); ok {
			location = trimFileName(file) + ":" + strconv.Itoa(line)
		} else {
			location = "???"
		}
		fmt.Fprintf(os.Stderr, "log: failed to parseFields, error=%v, location=%s\n", err, location)
		return l
	}
	m2 := make(map[string]interface{}, len(l.fields)+len(m))
	for k, v := range l.fields {
		m2[k] = v
	}
	for k, v := range m {
		m2[k] = v
	}
	return &logger{
		traceId:   l.traceId,
		fields:    m2,
		out:       l.out,
		formatter: l.formatter,
	}
}

var (
	_ErrNumberOfFieldsMustNotBeOdd = errors.New("the number of fields must not be odd")
	_ErrTypeOfFieldKeyMustBeString = errors.New("the type of field key must be string")
	_ErrFieldKeyMustNotBeEmpty     = errors.New("the field key must not be empty")
)

func parseFields(fields []interface{}) (map[string]interface{}, error) {
	if len(fields) == 0 {
		return nil, nil
	}
	if len(fields)&1 != 0 {
		return nil, _ErrNumberOfFieldsMustNotBeOdd
	}

	// 采用下面的实现可以避免边界检查.
	var (
		k  string
		ok bool
		m  = make(map[string]interface{}, len(fields)>>1)
	)
	for i, v := range fields {
		if i&1 == 0 { // key
			k, ok = v.(string)
			if !ok {
				return nil, _ErrTypeOfFieldKeyMustBeString
			}
			if k == "" {
				return nil, _ErrFieldKeyMustNotBeEmpty
			}
		} else { // value
			m[k] = v
		}
	}
	return m, nil
}
