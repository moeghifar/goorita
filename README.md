# Halaproxy
Simple & Minimalist http proxy written in Go

## Run
default run with default http port (33000)

```
go run main.go
```

custom http port

```
go run main.go -http-port=33001
```

run with dynamic nic binding

```
go run main.go -dynamic-nic-bind=enx:33001
```

## Version
Currently this proxy is in development version, we are going to release beta after some test and completion to features

## Features
- [x] Multi outgoing IP address
- [x] Multi port entry point
- [ ] HTTPS access
- [ ] Cached proxy
- [x] Run with flag configuration
- [ ] File based configuration
- [ ] Independent binary
- [ ] Serve with Fasthttp/Fiber

## Limitation
- Binding to Network Interface only work in Unix/Linux env since it needs `SO_BINDTODEVICE` 

## License
Apache License Version 2.0
