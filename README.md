# Meng
**M**ulti **E**ggress **N**etwork Proxy in **G**o.

A minimalist forward proxy written in Go.

Designated for multi eggress/outgoing network routing  

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
- [x] Access to HTTPS (with CONNECT method)
- [x] Flag based configuration
- [x] Better logging format
- [ ] Cached proxy
- [ ] File based configuration
- [ ] Independent binary
- [ ] Serve with Fasthttp/Fiber
- [ ] Debug log options from config 
- [ ] SOCKS5 and HTTPS entry support

## Limitation
- Binding to Network Interface only work in Unix/Linux env since it needs `SO_BINDTODEVICE` 

## References
- https://www.sobyte.net/post/2021-09/https-proxy-in-golang-in-less-than-100-lines-of-code/
- https://eli.thegreenplace.net/2022/go-and-proxy-servers-part-2-https-proxies/
- https://stackoverflow.com/a/27411300
- https://pkg.go.dev/syscall#BindToDevice

## License
Apache License Version 2.0
