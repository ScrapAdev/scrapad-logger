package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func (l *Logger) formatMessage(level string, message string) {
	switch level {
	case "info":
		l.formatInfoMessage(level, message)
	}
}

func (l *Logger) formatInfoMessage(level string, message string) {
	timestamp := aws.Time(time.Now().UTC()).UnixMilli()
	l.putLogEvent(timestamp, message, level)
}

func (l *Logger) formatErrorMessage(level string, message string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	timestamp := aws.Time(time.Now().UTC()).UnixMilli()
	logMessage := fmt.Sprintf(" [%d]:%s %s '%s:%d %s' - %s\n", timestamp, strings.Split(frame.Function, "/")[2], strings.ToUpper(level), frame.File, frame.Line, frame.Function[strings.LastIndex(frame.Function, ".")+1:], message)
	l.putLogEvent(timestamp, logMessage, level)
}
