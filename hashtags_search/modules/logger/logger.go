package logger

import (
	"encoding/json"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"

	logstash_logger "github.com/KaranJagtiani/go-logstash"
)

type LogEntry struct {
	Message    string `json:"message"`
	Level      string `json:"level"`
	TimeStamp  string `json:"timestamp"`
	Host       string `json:"host"`
	Method     string `json:"method"`
	Filename   string `json:"filename"`
	LineNumber string `json:"line_number"`
}

func CustomLogstashFormatter(level string, msg string) ([]byte, error) {

	logEntry := LogEntry{
		Message:    msg,
		Level:      level,
		TimeStamp:  time.Now().Format(time.RFC3339),
		Host:       GetHostName(),
		Method:     GetFuncName(1),
		Filename:   GetFileName(1),
		LineNumber: GetLineNumber(1),
	}

	return json.Marshal(logEntry)
}

func GetHostName() string {
	host, err := os.Hostname()

	if err != nil {
		return "unknown"
	}

	return host
}

func GetFuncName(i interface{}) string {

	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetFileName(skip int) string {
	_, file, _, ok := runtime.Caller(skip)

	if !ok {
		return "unknown"
	}

	return file
}

func GetLineNumber(skip int) string {

	_, _, line, ok := runtime.Caller(skip)

	if !ok {
		return "undefined"
	}

	return strconv.Itoa(line)
}

func LoggerIntialization() *logstash_logger.Logstash {

	logger := logstash_logger.Init("localhost", 5959, "udp", 5)

	return logger
}
