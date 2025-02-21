# ScrapAd Logger

ScrapAd logger module, to send logs from Go to AWS CloudWatch API

## Log Streams
1. Info
2. Error
3. Debug
4. Trace

## Message format
### Info / Debug / Error
1. Microservice name
2. Timestamp
3. Level
4. RequestID
5. File Path
6. Method name
7. Message

Example:
```json
{
    "service": "scrapad-logger",
    "timestamp": "2025-02-19T12:22:10Z",
    "level": "DEBUG",
    "request": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "message": "Test log message",
    "filePath": "/home/aalmaran/repos/scrapad-logger/logger/logger_test.go:101",
    "method": "captureStdout"
}
```
### Trace
1. Microservice name
2. Timestamp
3. Level
4. RequestID
5. Message

```json
{
    "service": "scrapad-logger",
    "timestamp": "2025-02-19T12:23:46Z",
    "level": "TRACE",
    "request": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "message": "Test log message with trace"
}
```
