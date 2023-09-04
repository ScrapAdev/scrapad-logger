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
```shell
scrapad-microservice:[2023-01-16T12:42:29.148+01:00] INFO 8f7642f8-19ba-4683-bf5d-dfb095fd0b55 'main/example/go.go:100' MethodTest - message example
```
### Trace
1. Microservice name
2. Timestamp
3. Level
4. RequestID
5. Message

```shell
scrapad-microservice:[2023-01-16T12:42:29.148+01:00] TRACE 8f7642f8-19ba-4683-bf5d-dfb095fd0b55 '200 Mozilla/5.0 (Macintosh) 192.168.0.1'
```
