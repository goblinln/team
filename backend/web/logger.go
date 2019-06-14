package web

import (
	"fmt"
	"os"
	"sync"
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

var (
	// LogName is prefix for all log files
	LogName = "log"
	// LogMaxLevel limits number of output messages.
	LogMaxLevel = LogLevelInfo
	// LogFileLimit limits size of per log file.
	LogFileLimit = 1024 * 1024 * 4
)

// FileLogger is a file-based log system.
type FileLogger struct {
	sync.Mutex

	writed     int
	createTime time.Time
	writer     *os.File
}

// Logger is runtime singleton instance of FileLogger.
var Logger = &FileLogger{
	Mutex:      sync.Mutex{},
	writed:     0,
	createTime: time.Now(),
	writer:     nil,
}

// Debug log output.
func (l *FileLogger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, "D", format, args...)
}

// Info writes normal log messages
func (l *FileLogger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, "I", format, args...)
}

// Warn writes messages as warning.
func (l *FileLogger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, "W", format, args...)
}

// Error writes messages as error.
func (l *FileLogger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, "E", format, args...)
}

// Fatal error output.
func (l *FileLogger) Fatal(format string, args ...interface{}) {
	l.log(LogLevelFatal, "F", format, args...)
	os.Exit(-1)
}

func (l *FileLogger) log(level int, prefix string, format string, args ...interface{}) {
	if level > LogMaxLevel {
		return
	}

	now := time.Now()
	msg := fmt.Sprintf(
		"[%04d-%02d-%02d %02d:%02d:%02d.%03d][%s]",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/(1000*1000),
		prefix) + fmt.Sprintf(format, args...)

	fmt.Println(msg)

	l.Lock()
	defer l.Unlock()

	if l.writer == nil || l.createTime.Day() != now.Day() || l.writed >= LogFileLimit {
		if l.writer != nil {
			l.writer.Close()
			l.writer = nil
		}

		dir := fmt.Sprintf("./logs/%04d%02d%02d", now.Year(), now.Month(), now.Day())
		if err := os.MkdirAll(dir, 777); err != nil {
			fmt.Printf("Failed to create log file! Errors: %s\n", err.Error())
			return
		}

		path := fmt.Sprintf("%s/%s_%02d%02d%02d.%03d.log", dir, LogName, now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/(1000*1000))
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			fmt.Printf("Failed to create log file! Errors: %s\n", err.Error())
			return
		}

		l.writed = 0
		l.createTime = now
		l.writer = file
	}

	n, err := l.writer.WriteString(msg + "\n")
	if err == nil {
		l.writed += n
	}
}
