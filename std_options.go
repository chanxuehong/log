package log

import (
	"fmt"
	"io"
	"sync"
)

func init() {
	SetFormatter(TextFormatter)
	SetOutput(concurrentStdout)
	SetLevel(DebugLevel)
}

var _stdOptions struct {
	sync.RWMutex
	options
}

func getStdOptions() (opts options) {
	_stdOptions.RLock()
	opts = _stdOptions.options
	_stdOptions.RUnlock()
	return
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	if formatter == nil {
		return
	}
	_stdOptions.Lock()
	_stdOptions.formatter = formatter
	_stdOptions.Unlock()
}

// SetOutput sets the standard logger output.
//  NOTE: output must be thread-safe, see ConcurrentWriter.
func SetOutput(output io.Writer) {
	if output == nil {
		return
	}
	_stdOptions.Lock()
	_stdOptions.output = output
	_stdOptions.Unlock()
}

func setLevel(level Level) {
	_stdOptions.Lock()
	_stdOptions.level = level
	_stdOptions.Unlock()
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
