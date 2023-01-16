package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

var (
	logGroupName = "arn:aws:logs:eu-west-1:690032394576:log-group:/scrapad-dev:*"
)

type Logger struct {
	logGroupName string
	svc          cloudwatchlogs.Client
}

func New() *Logger {
	options := cloudwatchlogs.Options{
		Region: *aws.String(os.Getenv("AWS_REGION")),
	}
	svc := cloudwatchlogs.New(options)
	return &Logger{
		logGroupName: os.Getenv("AWS_GROUP_NAME"),
		svc:          *svc,
	}
}

// [2023-01-12 11:07:24]:scrapad-query(PID) WARN 'main/example/go.go:100' MethodTest - message example
// [2023-01-12 11:07:24] '/scrapad-dev/graph/schema.resolvers.go:461 github.com/ScrapAdev/scrapad-query/graph.(*queryResolver).GetDashboard' - NO USER

func (log *Logger) Debug(txt string) {
	log.formatMessage("debug", txt)
}

func (log *Logger) Error(txt string) {
	log.formatMessage("error", txt)
}

func (log *Logger) Info(txt string) {
	log.formatMessage("info", txt)
}

func (log *Logger) formatMessage(level string, message string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	timestamp := aws.Time(time.Now().UTC()).UnixMilli()
	logMessage := fmt.Sprintf(" [%d]:%s %s '%s:%d %s' - %s\n", timestamp, strings.Split(frame.Function, "/")[2], strings.ToUpper(level), frame.File, frame.Line, frame.Function[strings.LastIndex(frame.Function, ".")+1:], message)
	log.putLogEvent(timestamp, logMessage, level)
}

func (log *Logger) putLogEvent(timestamp int64, msg string, level string) error {

	inputLog := types.InputLogEvent{
		Message:   aws.String(msg),
		Timestamp: &timestamp,
	}
	inputLogEvent := []types.InputLogEvent{inputLog}
	params := cloudwatchlogs.PutLogEventsInput{
		LogEvents:     inputLogEvent,
		LogGroupName:  aws.String(log.logGroupName),
		LogStreamName: aws.String(level),
	}

	_, err := log.svc.PutLogEvents(nil, &params)

	if err != nil {
		fmt.Println("Failed to send log request", err)
		return err
	}
	return err
}
