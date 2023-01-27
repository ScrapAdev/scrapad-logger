package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func (l *Logger) Debug(txt string) {
	l.formatMessage("debug", txt)
}

func (l *Logger) Error(txt string) {
	l.initframes()
	l.formatMessage("error", txt)
}

func (l *Logger) Info(txt string) {
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
	logMessage := fmt.Sprintf("%s:[%s] %s %d '%s'\n", "TODO", timestamp.String(), strings.ToUpper(level), 0, message)
	fmt.Println(logMessage)
	l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
}

func (l *Logger) formatOtherMessage(level string, message string) {
	timestamp := aws.Time(time.Now().UTC())
	service := "no-service"
	repo := strings.Split(l.frame.Function, "/")
	if len(repo) > 1 {
		service = strings.Split(repo[2], ".")[0]
	}
	logMessage := fmt.Sprintf("%s:[%s] %s '%s:%d %s' - %s\n", service, timestamp.String(), strings.ToUpper(level), l.frame.File, l.frame.Line, l.frame.Function[strings.LastIndex(l.frame.Function, ".")+1:], message)
	fmt.Println(logMessage)
	l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
}
