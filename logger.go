package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type Logger struct {
	logGroupName string
	svc          cloudwatchlogs.Client
}

// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Println("Error loading .env file")
// 	}
// 	l := New()
// 	l.Error("ERROR")
// }

func New() *Logger {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		fmt.Println(cfg)
	}

	svc := cloudwatchlogs.NewFromConfig(cfg)
	return &Logger{
		logGroupName: os.Getenv("AWS_GROUP_NAME"),
		svc:          *svc,
	}
}

func (l *Logger) putLogEvent(timestamp int64, msg string, level string) error {

	inputLog := types.InputLogEvent{
		Message:   aws.String(msg),
		Timestamp: &timestamp,
	}
	inputLogEvent := []types.InputLogEvent{inputLog}
	params := cloudwatchlogs.PutLogEventsInput{
		LogEvents:     inputLogEvent,
		LogGroupName:  aws.String(l.logGroupName),
		LogStreamName: aws.String(level),
	}
	_, err := l.svc.PutLogEvents(context.Background(), &params)
	if err != nil {
		fmt.Println("Failed to send log request", err)
		return err
	}
	return err
}

func (l *Logger) Debug(txt string) {
	l.formatMessage("debug", txt)
}

func (l *Logger) Error(txt string) {
	l.formatMessage("error", txt)
}

func (l *Logger) Info(txt string) {
	l.formatMessage("info", txt)
}

func (l *Logger) Trace(txt string) {
	l.formatMessage("trace", txt)
}

func (l *Logger) formatMessage(level string, message string) {
	fmt.Println("Logging...", message)
	if level == "trace" {
		l.formatTraceMessage(level, message)
	}
	l.formatOtherMessage(level, message)
}

func (l *Logger) formatTraceMessage(level string, message string) {
	timestamp := aws.Time(time.Now().UTC())
	logMessage := fmt.Sprintf("%s:[%s] %s %d '%s'\n", "TODO", timestamp.String(), strings.ToUpper(level), 0, message)
	l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
}

func (l *Logger) formatOtherMessage(level string, message string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	timestamp := aws.Time(time.Now().UTC())
	logMessage := fmt.Sprintf("%s:[%s] %s '%s:%d %s' - %s\n", strings.Split(frame.Function, "/")[0], timestamp.String(), strings.ToUpper(level), frame.File, frame.Line, frame.Function[strings.LastIndex(frame.Function, ".")+1:], message)
	fmt.Println(logMessage)
	err := l.putLogEvent(timestamp.UnixMilli(), logMessage, level)
	if err != nil {
		fmt.Println(err.Error())
	}
}
