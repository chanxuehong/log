package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chanxuehong/log/trace"
)

type loggerContextKey struct{}

func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return New(WithTraceId(trace.NewTraceId()))
	}
	v, ok := ctx.Value(loggerContextKey{}).(Logger)
	if ok && v != nil {
		return v
	}
	traceId := trace.FromContext(ctx)
	return New(WithTraceId(traceId))
}

func FromRequest(req *http.Request) Logger {
	if req == nil {
		return New(WithTraceId(trace.NewTraceId()))
	}
	v, ok := req.Context().Value(loggerContextKey{}).(Logger)
	if ok && v != nil {
		return v
	}
	traceId := trace.FromRequest(req)
	return New(WithTraceId(traceId))
}

func FromHeader(header http.Header) Logger {
	traceId := trace.FromHeader(header)
	return New(WithTraceId(traceId))
}

func NewContext(ctx context.Context, logger Logger) context.Context {
	if logger == nil {
		return ctx
	}
	if ctx == nil {
		return context.WithValue(context.Background(), loggerContextKey{}, logger)
	}
	if v, ok := ctx.Value(loggerContextKey{}).(Logger); ok && v == logger {
		return ctx
	}
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

type Logger interface {
	// Fatal logs a message at FatalLevel.
	//
	// Unlike other golang log libraries (for example, the golang standard log library),
	// Fatal just logs a message and does not call os.Exit, so you need to explicitly call os.Exit if necessary.
	//
	// For fields, the following conditions must be satisfied
	//  1. the len(fields) must be an even number, that is to say len(fields)%2==0
	//  2. the even index element of fields must be non-empty string
	Fatal(msg string, fields ...interface{})

	// Error logs a message at ErrorLevel.
	// The requirements for fields can see the comments of Fatal.
	Error(msg string, fields ...interface{})

	// Warn logs a message at WarnLevel.
	// The requirements for fields can see the comments of Fatal.
	Warn(msg string, fields ...interface{})

	// Info logs a message at InfoLevel.
	// The requirements for fields can see the comments of Fatal.
	Info(msg string, fields ...interface{})

	// Debug logs a message at DebugLevel.
	// The requirements for fields can see the comments of Fatal.
	Debug(msg string, fields ...interface{})

	// Output logs a message at specified level.
	//
	// For level==FatalLevel, unlike other golang log libraries (for example, the golang standard log library),
	// Output just logs a message and does not call os.Exit, so you need to explicitly call os.Exit if necessary.
	//
	// The requirements for fields can see the comments of Fatal.
	Output(calldepth int, level Level, msg string, fields ...interface{})

	// WithField creates a new Logger from the current Logger and adds a field to it.
	WithField(key string, value interface{}) Logger

	// WithFields creates a new Logger from the current Logger and adds multiple fields to it.
	// The requirements for fields can see the comments of Fatal.
	WithFields(fields ...interface{}) Logger

	// SetFormatter sets the logger formatter.
	SetFormatter(formatter Formatter)

	// SetOutput sets the logger output.
	SetOutput(output io.Writer)

	// SetLevel sets the logger level.
	SetLevel(level Level) error

	// SetLevelString sets the logger level.
	SetLevelString(str string) error
}

type Option func(*options)

func WithTraceId(traceId string) Option {
	return func(o *options) {
		o.traceId = traceId
	}
}

// WithOutput sets the logger output.
func WithOutput(output io.Writer) Option {
	return func(o *options) {
		if output == nil {
			return
		}
		o.output = output
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		if formatter == nil {
			return
		}
		o.formatter = formatter
	}
}

func WithLevel(level Level) Option {
	return func(o *options) {
		if !isValidLevel(level) {
			return
		}
		o.level = level
	}
}

func WithLevelString(str string) Option {
	return func(o *options) {
		level, ok := parseLevelString(str)
		if !ok {
			return
		}
		o.level = level
	}
}

func New(opts ...Option) Logger { return _New(opts) }

func _New(opts []Option) *logger {
	l := &logger{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&l.options)
	}
	if l.options.formatter == nil {
		l.options.formatter = TextFormatter
	}
	if l.options.output == nil {
		l.options.output = os.Stdout
	}
	if l.options.level == InvalidLevel {
		l.options.level = DebugLevel
	}
	return l
}

type logger struct {
	mu      sync.Mutex // ensures atomic writes; protects the following options field
	options options
	fields  map[string]interface{}
}

type options struct {
	traceId   string
	formatter Formatter
	output    io.Writer
	level     Level
}

type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

type Entry struct {
	Location string // function(file:line)
	Time     time.Time
	Level    Level
	TraceId  string
	Message  string
	Fields   map[string]interface{}
	Buffer   *bytes.Buffer
}

func (l *logger) SetFormatter(formatter Formatter) {
	if formatter == nil {
		return
	}
	l.mu.Lock()
	l.options.formatter = formatter
	l.mu.Unlock()
}
func (l *logger) SetOutput(output io.Writer) {
	if output == nil {
		return
	}
	l.mu.Lock()
	l.options.output = output
	l.mu.Unlock()
}
func (l *logger) SetLevel(level Level) error {
	if !isValidLevel(level) {
		return fmt.Errorf("invalid level: %d", level)
	}
	l.setLevel(level)
	return nil
}
func (l *logger) SetLevelString(str string) error {
	level, ok := parseLevelString(str)
	if !ok {
		return fmt.Errorf("invalid level string: %q", str)
	}
	l.setLevel(level)
	return nil
}
func (l *logger) setLevel(level Level) {
	l.mu.Lock()
	l.options.level = level
	l.mu.Unlock()
}
func (l *logger) getOptions() (opts options) {
	l.mu.Lock()
	opts = l.options
	l.mu.Unlock()
	return
}

func (l *logger) Fatal(msg string, fields ...interface{}) {
	l.output(1, FatalLevel, msg, fields)
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

func (l *logger) Output(calldepth int, level Level, msg string, fields ...interface{}) {
	if !isValidLevel(level) {
		return
	}
	if calldepth < 0 {
		calldepth = 0
	}
	l.output(calldepth+1, level, msg, fields)
}

func (l *logger) output(calldepth int, level Level, msg string, fields []interface{}) {
	opts := l.getOptions()
	if !isLevelEnabled(level, opts.level) {
		return
	}

	var location string
	if pc, file, line, ok := runtime.Caller(calldepth + 1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			location = trimFuncName(fn.Name()) + "(" + trimFileName(file) + ":" + strconv.Itoa(line) + ")"
		} else {
			location = trimFileName(file) + ":" + strconv.Itoa(line)
		}
	} else {
		location = "???"
	}

	var m map[string]interface{}
	if len(fields) == 0 {
		m = l.fields
	} else {
		m2, err := parseFields(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "log: failed to parse fields, error=%v, location=%s\n", err, location)
		}
		if len(m2) == 0 {
			m = l.fields
		} else {
			if len(l.fields) == 0 {
				m = m2
			} else {
				m = make(map[string]interface{}, len(l.fields)+len(m2))
				for k, v := range l.fields {
					m[k] = v
				}
				for k, v := range m2 {
					m[k] = v
				}
			}
		}
	}

	buffer := _bufferPool.Get().(*bytes.Buffer)
	defer _bufferPool.Put(buffer)
	buffer.Reset()

	data, err := opts.formatter.Format(&Entry{
		Location: location,
		Time:     time.Now(),
		Level:    level,
		TraceId:  opts.traceId,
		Message:  msg,
		Fields:   m,
		Buffer:   buffer,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to format Entry, error=%v, location=%s\n", err, location)
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err = opts.output.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to write to log, error=%v, location=%s\n", err, location)
		return
	}
}

var _bufferPool = sync.Pool{
	New: newBuffer,
}

func newBuffer() interface{} {
	return bytes.NewBuffer(make([]byte, 0, 16<<10))
}

func trimFuncName(name string) string {
	return path.Base(name)
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
	if key == "" {
		return l
	}
	m := make(map[string]interface{}, len(l.fields)+1)
	for k, v := range l.fields {
		m[k] = v
	}
	m[key] = value
	return &logger{
		options: l.getOptions(),
		fields:  m,
	}
}
func (l *logger) WithFields(fields ...interface{}) Logger {
	if len(fields) == 0 {
		return l
	}
	m, err := parseFields(fields)
	if err != nil {
		var location string
		if pc, file, line, ok := runtime.Caller(1); ok {
			if fn := runtime.FuncForPC(pc); fn != nil {
				location = trimFuncName(fn.Name()) + "(" + trimFileName(file) + ":" + strconv.Itoa(line) + ")"
			} else {
				location = trimFileName(file) + ":" + strconv.Itoa(line)
			}
		} else {
			location = "???"
		}
		fmt.Fprintf(os.Stderr, "log: failed to parse fields, error=%v, location=%s\n", err, location)
	}
	if len(m) == 0 {
		return l
	}
	if len(l.fields) == 0 {
		return &logger{
			options: l.getOptions(),
			fields:  m,
		}
	}
	m2 := make(map[string]interface{}, len(l.fields)+len(m))
	for k, v := range l.fields {
		m2[k] = v
	}
	for k, v := range m {
		m2[k] = v
	}
	return &logger{
		options: l.getOptions(),
		fields:  m2,
	}
}

var (
	_ErrNumberOfFieldsMustNotBeOdd error = errors.New("the number of fields must not be odd")
	_ErrTypeOfFieldKeyMustBeString error = errors.New("the type of field key must be string")
	_ErrFieldKeyMustNotBeEmpty     error = errors.New("the field key must not be empty")
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
				return m, _ErrTypeOfFieldKeyMustBeString
			}
			if k == "" {
				return m, _ErrFieldKeyMustNotBeEmpty
			}
		} else { // value
			m[k] = v
		}
	}
	return m, nil
}

var _ trace.Tracer = (*logger)(nil)

func (l *logger) TraceId() string {
	return l.getOptions().traceId
}
