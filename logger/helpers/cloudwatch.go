package helpers

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type AWSCloudwatch struct {
	client       *cloudwatchlogs.Client
	logGroupName string
}

func NewAWSCloudwatch() (*AWSCloudwatch, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, fmt.Errorf("error creating AWS config: %v", err)
	}

	svc := cloudwatchlogs.NewFromConfig(cfg)
	return &AWSCloudwatch{
		client:       svc,
		logGroupName: os.Getenv("AWS_GROUP_NAME"),
	}, nil
}

func (cloudwatch *AWSCloudwatch) PutLogEvent(timestamp int64, msg string, level string) error {
	inputLog := types.InputLogEvent{
		Message:   aws.String(msg),
		Timestamp: &timestamp,
	}
	inputLogEvent := []types.InputLogEvent{inputLog}
	params := cloudwatchlogs.PutLogEventsInput{
		LogEvents:     inputLogEvent,
		LogGroupName:  aws.String(cloudwatch.logGroupName),
		LogStreamName: aws.String(level),
	}
	_, err := cloudwatch.client.PutLogEvents(context.Background(), &params)
	if err != nil {
		return fmt.Errorf("failed to send log request: %v", err)
	}
	return nil
}
