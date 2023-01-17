package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

func (l *Logger) formatMessage(level string, message string) {
	if level == "trace" {
		l.formatTraceMessage(level, message)
	}
	l.formatOtherMessage(level, message)
}

func (l *Logger) formatOtherMessage(level string, message string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	timestamp := aws.Time(time.Now().UTC())
	logMessage := fmt.Sprintf("%s:[%s] %s '%s:%d %s' - %s\n", strings.Split(frame.Function, "/")[2], timestamp.String(), strings.ToUpper(level), frame.File, frame.Line, frame.Function[strings.LastIndex(frame.Function, ".")+1:], message)
	l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
}

func (l *Logger) formatTraceMessage(level string, message string) {
	timestamp := aws.Time(time.Now().UTC())
	logMessage := fmt.Sprintf("%s:[%s] %s %d '%s'\n", "TODO", timestamp.String(), strings.ToUpper(level), 0, message)
	l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
}
