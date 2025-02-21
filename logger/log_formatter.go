package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

var (
	CALLER_SKIP  = 3
	CALLER_DEPTH = 32
)

func (l *Logger) Debug(txt string) {
	level := "debug"
	if l.isLogLevelActivated(level) {
		l.formatLogMessage(level, txt)
	}
}

func (l *Logger) Error(txt string) {
	level := "error"
	if l.isLogLevelActivated(level) {
		l.formatLogMessage(level, txt)
	}
}

func (l *Logger) Info(txt string) {
	level := "info"
	if l.isLogLevelActivated(level) {
		l.formatLogMessage(level, txt)
	}
}

func (l *Logger) Trace(txt string) {
	level := "trace"
	if l.isLogLevelActivated(level) {
		l.formatLogMessage(level, txt)
	}
}

type logMessage struct {
	Service   string  `json:"service"`
	Timestamp string  `json:"timestamp"`
	Level     string  `json:"level"`
	Request   string  `json:"request"`
	Message   string  `json:"message"`
	FilePath  *string `json:"filePath,omitempty"`
	Method    *string `json:"method,omitempty"`
}

func (l *Logger) formatLogMessage(level, message string) {
	l.initFrames()
	timestamp, service := l.parseMessageInfo()

	logMessage := logMessage{
		Service:   service,
		Timestamp: timestamp.Format(time.RFC3339),
		Level:     strings.ToUpper(level),
		Request:   l.uuidRequest.String(),
		Message:   message,
	}
	if level != "trace" {
		filePath := fmt.Sprintf("%s:%d", l.frame.File, l.frame.Line)
		method := l.frame.Function[strings.LastIndex(l.frame.Function, ".")+1:]

		logMessage.FilePath = &filePath
		logMessage.Method = &method
	}

	logMessageBytes, _ := json.Marshal(logMessage)
	fmt.Println(string(logMessageBytes))
	l.cloudwatch.PutLogEvent(timestamp.UnixMilli(), string(logMessageBytes), level)
}

func (l *Logger) isLogLevelActivated(level string) bool {
	var isActivated bool
	for _, v := range strings.Split(os.Getenv("LOG_LEVELS"), ",") {
		if v == level {
			isActivated = true
		}
	}
	return isActivated
}

func (l *Logger) parseMessageInfo() (*time.Time, string) {
	timestamp := aws.Time(time.Now().UTC())
	service := "no-service"
	repo := strings.Split(l.frame.Function, "/")
	if len(repo) > 1 {
		service = strings.Split(repo[2], ".")[0]
	}
	return timestamp, service
}

func (l *Logger) printCallStack() {
	debug.PrintStack()
	callers := make([]uintptr, CALLER_DEPTH)
	n := runtime.Callers(CALLER_SKIP, callers)
	frames := runtime.CallersFrames(callers[:n])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		fmt.Printf("+ %s:%d %s\n", frame.File, frame.Line, frame.Function)
	}
}

func (l *Logger) initFrames() {
	pc := make([]uintptr, CALLER_DEPTH)
	n := runtime.Callers(CALLER_SKIP, pc)
	l.frames = runtime.CallersFrames(pc[:n])
	l.frames.Next()
	l.frame, _ = l.frames.Next()
	l.frames = runtime.CallersFrames(pc[:n])
}
