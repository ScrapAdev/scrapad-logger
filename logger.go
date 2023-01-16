package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/joho/godotenv"
)

type Logger struct {
	logGroupName string
	svc          cloudwatchlogs.Client
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	l := New()
	l.Info("Segundo test")
}

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

func (l *Logger) Debug(txt string) {
	l.formatMessage("debug", txt)
}

func (l *Logger) Error(txt string) {
	l.formatMessage("error", txt)
}

func (l *Logger) Info(txt string) {
	l.formatMessage("info", txt)
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
