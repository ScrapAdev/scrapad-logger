package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
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
	logMessage := "Test log message"

	captured := captureStdout(func() {
		logger.formatLogMessage(logLevel, logMessage)
	})

	assert.Contains(t, captured, strings.ToUpper(logLevel))
	assert.Contains(t, captured, uuidStr)
	assert.Contains(t, captured, logMessage)
}

func TestLogger_FormatTraceMessage(t *testing.T) {
	uuidStr := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	logger, err := NewWithID(uuidStr)
	assert.NoError(t, err)

	logLevel := "trace"
	logMessage := "Test log message with trace"

	captured := captureStdout(func() {
		logger.formatLogMessage(logLevel, logMessage)
	})

	assert.Contains(t, captured, strings.ToUpper(logLevel))
	assert.Contains(t, captured, uuidStr)
	assert.Contains(t, captured, logMessage)
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

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String()
}
