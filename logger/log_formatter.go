package logger

import (
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

func (l *Logger) formatLogMessage(level, message string) string {
	l.initFrames()
	timestamp, service := l.parseMessageInfo()

	var logMessage string
	if level == "trace" {
		logMessage = fmt.Sprintf("%s:[%s] %s %s '%s'\n", service, timestamp.Format("2006-01-02 15:04:05"), strings.ToUpper(level), l.uuidRequest.String(), message)
	} else {
		logMessage = fmt.Sprintf("%s:[%s] %s %s '%s:%d %s' - %s\n", service, timestamp.Format("2006-01-02 15:04:05"), strings.ToUpper(level), l.uuidRequest.String(), l.frame.File, l.frame.Line, l.frame.Function[strings.LastIndex(l.frame.Function, ".")+1:], message)
	}

	fmt.Println(logMessage)
	l.cloudwatch.PutLogEvent(timestamp.UnixMilli(), logMessage, level)

	return logMessage
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
