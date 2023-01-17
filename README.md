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
4. File Path
5. Method name
6. Message

Example:
```shell
scrapad-microservice:[2023-01-16T12:42:29.148+01:00] INFO 'main/example/go.go:100' MethodTest - message example
```
### Trace
1. Microservice
2. Timestamp
3. Level
4. HTTP Code
5. URL
6. Browser Info
7. User IP

```shell
scrapad-microservice:[2023-01-16T12:42:29.148+01:00] TRACE 200 'Mozilla/5.0 (Macintosh)' "192.168.0.1"
```
