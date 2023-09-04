package logger

import (
	"runtime"

	"github.com/ScrapAdev/scrapad-logger/logger/helpers"
	"github.com/google/uuid"
)

type Logger struct {
	cloudwatch  helpers.AWSCloudwatch
	frames      *runtime.Frames
	frame       runtime.Frame
	uuidRequest uuid.UUID
}

func New(uuidRequest string) (*Logger, error) {
	awsCloudwatch, err := helpers.NewAWSCloudwatch()
	if err != nil {
		return nil, err
	}
	uuid, err := uuid.Parse(uuidRequest)
	if err != nil {
		return nil, err
	}
	return &Logger{
		cloudwatch:  *awsCloudwatch,
		uuidRequest: uuid,
	}, nil
}
