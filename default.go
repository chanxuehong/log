package log

import (
	"io"
	"sync"
)

var (
	_formatterMutex sync.RWMutex
	_formatter      Formatter

	_outputMutex sync.RWMutex
	_output      io.Writer
)

func init() {
	SetFormatter(TextFormatter)
	SetOutput(concurrentStdout)
}

func getFormatter() Formatter {
	_formatterMutex.RLock()
	formatter := _formatter
	_formatterMutex.RUnlock()
	return formatter
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	if formatter == nil {
		return
	}
	_formatterMutex.Lock()
	_formatter = formatter
	_formatterMutex.Unlock()
}

func getOutput() io.Writer {
	_outputMutex.RLock()
	output := _output
	_outputMutex.RUnlock()
	return output
}

// SetOutput sets the standard logger output.
//  NOTE: output must be thread-safe, see ConcurrentWriter.
func SetOutput(output io.Writer) {
	if output == nil {
		return
	}
	_outputMutex.Lock()
	_output = output
	_outputMutex.Unlock()
}
