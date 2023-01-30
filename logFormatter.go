package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func (l *Logger) Debug(txt string) {
	l.initFrames()
	l.formatMessage("debug", txt)
}

func (l *Logger) Error(txt string) {
	l.initFrames()
	l.formatMessage("error", txt)
}

func (l *Logger) Info(txt string) {
	l.initFrames()
	l.formatMessage("info", txt)
}

func (l *Logger) Trace(txt string) {
	l.formatMessage("trace", txt)
}

func (l *Logger) formatMessage(level string, message string) {
	if level == "trace" {
		l.formatTraceMessage(level, message)
	}
	l.formatOtherMessage(level, message)
}

func (l *Logger) formatTraceMessage(level string, message string) {
	timestamp := aws.Time(time.Now().UTC())
	logMessage := fmt.Sprintf("%s:[%s] %s %d '%s'\n", "TODO", timestamp.Format("2006-01-02 15:04:05"), strings.ToUpper(level), 0, message)
	fmt.Println(logMessage)
	l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
}

func (l *Logger) formatOtherMessage(level string, message string) {
	timestamp, service := l.parseMessageInfo()
	logMessage := fmt.Sprintf("%s:[%s] %s '%s:%d %s' - %s\n", service, timestamp.Format("2006-01-02 15:04:05"), strings.ToUpper(level), l.frame.File, l.frame.Line, l.frame.Function[strings.LastIndex(l.frame.Function, ".")+1:], message)
	fmt.Println(logMessage)
	l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
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

func (l *Logger) initFrames() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	l.frames = runtime.CallersFrames(pc[:n])
	l.frames.Next()
	l.frame, _ = l.frames.Next()
	l.frames = runtime.CallersFrames(pc[:n])
}

// func (l *Logger) callStack() string {
// 	frames := l.frames
// 	callStack := []string{}
// 	for {
// 		frame, more := frames.Next()
// 		if !more {
// 			break
// 		}
// 		callStack = append(callStack, fmt.Sprintf("%s:%d", frame.Function, frame.Line))
// 	}
// 	stack := fmt.Sprint(callStack)
// 	return stack
// }
