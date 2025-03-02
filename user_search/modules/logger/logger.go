package logger

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"time"

	logstash "github.com/KaranJagtiani/go-logstash"
)

// LogEntry represents a structured log entry.
type LogEntry struct {
	Message    string `json:"message"`
	Level      string `json:"level"`
	TimeStamp  string `json:"timestamp"`
	Host       string `json:"host"`
	Method     string `json:"method"`
	Filename   string `json:"filename"`
	LineNumber string `json:"line_number"`
}

// FormatLogEntry creates a structured log entry in JSON format.
func FormatLogEntry(level, msg string) ([]byte, error) {
	logEntry := LogEntry{
		Message:    msg,
		Level:      level,
		TimeStamp:  time.Now().Format(time.RFC3339),
		Host:       getHostName(),
		Method:     getCallerFunction(2),
		Filename:   getCallerFile(2),
		LineNumber: getCallerLine(2),
	}
	return json.Marshal(logEntry)
}

// getHostName retrieves the system hostname.
func getHostName() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

// getCallerFunction returns the name of the function that called the logger.
func getCallerFunction(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return runtime.FuncForPC(pc).Name()
}

// getCallerFile returns the filename where the logger was called.
func getCallerFile(skip int) string {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return file
}

// getCallerLine returns the line number where the logger was called.
func getCallerLine(skip int) string {
	_, _, line, ok := runtime.Caller(skip)
	if !ok {
		return "undefined"
	}
	return strconv.Itoa(line)
}

// InitializeLogger sets up a Logstash logger instance.
func InitializeLogger() *logstash.Logstash {
	return logstash.Init("localhost", 5959, "udp", 5)
}