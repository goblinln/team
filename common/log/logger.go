package log

import (
	"fmt"
	"os"
	"time"
)

// Limit for log outputs.
const (
	// Only output fatal message.
	LogLevelFatal = 0
	// Both error and fatal messages will be outputed.
	LogLevelError = 1
	// Warn/Error/Fatal
	LogLevelWarn = 2
	// Debug message will not be outputed.
	LogLevelInfo = 3
	// All messages will be found in log.
	LogLevelDebug = 4
)

// MaxLogLevel controls output messages kinds.
var MaxLogLevel = LogLevelInfo

// Debug log output.
func Debug(format string, args ...interface{}) {
	log(LogLevelDebug, "D", format, args...)
}

// Info writes normal log messages
func Info(format string, args ...interface{}) {
	log(LogLevelInfo, "I", format, args...)
}

// Warn writes messages as warning.
func Warn(format string, args ...interface{}) {
	log(LogLevelWarn, "W", format, args...)
}

// Error writes messages as error.
func Error(format string, args ...interface{}) {
	log(LogLevelError, "E", format, args...)
}

// Fatal error output.
func Fatal(format string, args ...interface{}) {
	log(LogLevelFatal, "F", format, args...)
	os.Exit(-1)
}

func log(level int, prefix string, format string, args ...interface{}) {
	if level > MaxLogLevel {
		return
	}

	now := time.Now()
	msg := fmt.Sprintf(
		"[%04d-%02d-%02d %02d:%02d:%02d.%03d][%s]",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/(1000*1000),
		prefix) + fmt.Sprintf(format, args...)

	fmt.Println(msg)
}
