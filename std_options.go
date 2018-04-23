package log

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

var _stdOptions struct {
	formatterMutex sync.RWMutex
	formatter      Formatter

	outputMutex sync.RWMutex
	output      io.Writer

	level uint64 // Level
}

func init() {
	SetFormatter(TextFormatter)
	SetOutput(concurrentStdout)
	SetLevel(DebugLevel)
}

func getFormatter() Formatter {
	_stdOptions.formatterMutex.RLock()
	formatter := _stdOptions.formatter
	_stdOptions.formatterMutex.RUnlock()
	return formatter
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	if formatter == nil {
		return
	}
	_stdOptions.formatterMutex.Lock()
	_stdOptions.formatter = formatter
	_stdOptions.formatterMutex.Unlock()
}

func getOutput() io.Writer {
	_stdOptions.outputMutex.RLock()
	output := _stdOptions.output
	_stdOptions.outputMutex.RUnlock()
	return output
}

// SetOutput sets the standard logger output.
//  NOTE: output must be thread-safe, see ConcurrentWriter.
func SetOutput(output io.Writer) {
	if output == nil {
		return
	}
	_stdOptions.outputMutex.Lock()
	_stdOptions.output = output
	_stdOptions.outputMutex.Unlock()
}

func getLevel() Level {
	return Level(atomic.LoadUint64(&_stdOptions.level))
}

func setLevel(level Level) {
	atomic.StoreUint64(&_stdOptions.level, uint64(level))
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) error {
	if !isValidLevel(level) {
		return fmt.Errorf("invalid level: %d", level)
	}
	setLevel(level)
	return nil
}

// SetLevelString sets the standard logger level.
func SetLevelString(str string) error {
	level, ok := parseLevelString(str)
	if !ok {
		return fmt.Errorf("invalid level string: %q", str)
	}
	setLevel(level)
	return nil
}
