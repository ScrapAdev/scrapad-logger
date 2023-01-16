# Scrapad Logger

## Message format
### Log | Debug | Error
1. IP máquina
2. Microservicio
3. Timestamp
4. Level
5. File Path
6. Method name
7. Message

Example:
```shell
scrapad-query:(192.168.0.1) [2023-01-16T12:42:29.148+01:00] INFO 'main/example/go.go:100' MethodTest - message example
```
### Trace
1. IP máquina
2. Microservicio
3. Timestamp
4. Level
5. HTTP Code
6. URL
7. Browser Info
8. User IP

```shell
scrapad-query:(192.168.0.1) [2023-01-16T12:42:29.148+01:00] TRACE 200 'Mozilla/5.0 (Macintosh)' "83.122.11.0"
```
