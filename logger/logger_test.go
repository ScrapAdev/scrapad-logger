package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	uuidStr := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	logger, err := NewWithID(uuidStr)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	invalidUUIDStr := "invalid-uuid"
	_, err = NewWithID(invalidUUIDStr)
	assert.Error(t, err)
}

func TestLogger_FormatMessage(t *testing.T) {
	uuidStr := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	logger, err := NewWithID(uuidStr)
	assert.NoError(t, err)

	logLevel := "debug"
	message := "Test log message"

	captured := captureStdout(func() {
		logger.formatLogMessage(logLevel, message)
	})

	var gotMessage logMessage
	_ = json.Unmarshal([]byte(captured), &gotMessage)

	assert.Equal(t, gotMessage.Service, "scrapad-logger")
	assert.Equal(t, gotMessage.Level, "DEBUG")
	assert.Equal(t, gotMessage.Request, "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	assert.Equal(t, gotMessage.Message, "Test log message")

	assert.Equal(t, *gotMessage.Method, "captureStdout")
	assert.Contains(t, *gotMessage.FilePath, "scrapad-logger/logger/logger_test.go")
}

func TestLogger_FormatTraceMessage(t *testing.T) {
	uuidStr := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	logger, err := NewWithID(uuidStr)
	assert.NoError(t, err)

	logLevel := "trace"
	message := "Test log message with trace"

	captured := captureStdout(func() {
		logger.formatLogMessage(logLevel, message)
	})

	var gotMessage logMessage
	_ = json.Unmarshal([]byte(captured), &gotMessage)

	assert.Equal(t, gotMessage.Service, "scrapad-logger")
	assert.Equal(t, gotMessage.Level, "TRACE")
	assert.Equal(t, gotMessage.Request, "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	assert.Equal(t, gotMessage.Message, "Test log message with trace")
	assert.Nil(t, gotMessage.FilePath)
	assert.Nil(t, gotMessage.Method)
}

func TestLogger_InitFrames(t *testing.T) {
	uuidStr := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	logger, err := NewWithID(uuidStr)
	assert.NoError(t, err)

	logger.initFrames()

	assert.NotNil(t, logger.frames)
	assert.NotNil(t, logger.frame)
	assert.NotNil(t, logger.frame.File)
	assert.NotNil(t, logger.frame.Line)
	assert.NotNil(t, logger.frame.Function)
}

func TestLogger_AllMethods(t *testing.T) {
	uuidStr := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	logger, err := NewWithID(uuidStr)
	assert.NoError(t, err)

	logger.Debug("This is a debug message")
	logger.Error("This is an error message")
	logger.Info("This is an info message")
	logger.Trace("This is a trace message")
}

func TestLogger_ParseMessageInfo_ShortFunctionPaths(t *testing.T) {
	logger, err := New()
	assert.NoError(t, err)

	tests := []struct {
		name            string
		functionName    string
		expectedService string
	}{
		{
			name:            "empty function",
			functionName:    "",
			expectedService: "no-service",
		},
		{
			name:            "single part function (no slash)",
			functionName:    "main.main",
			expectedService: "no-service",
		},
		{
			name:            "two parts function (single slash)",
			functionName:    "command-line-arguments/server.(*grpcServer).GetTransactionPrices",
			expectedService: "no-service",
		},
		{
			name:            "standard three parts function",
			functionName:    "github.com/ScrapAdev/scrapad-proxy/internal/handler.handleV2Request",
			expectedService: "scrapad-proxy",
		},
		{
			name:            "standard payments function",
			functionName:    "github.com/ScrapAdev/scrapad-payment-service/app.(*grpcServer).GetTransactionPrices",
			expectedService: "scrapad-payment-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.frame.Function = tt.functionName
			_, service := logger.parseMessageInfo()
			assert.Equal(t, tt.expectedService, service)
		})
	}
}

func TestLogger_Concurrency(t *testing.T) {
	logger, err := New()
	assert.NoError(t, err)

	const goroutines = 50
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			logger.formatLogMessage("info", "concurrent test log")
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}
