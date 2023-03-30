# goorita
Simple http proxy written in Go

## limitation
- Binding to Network Interface only work in Unix/Linux env since it needs `SO_BINDTODEVICE` 

## features
- [x] Multi outgoing IP address
- [x] Multi port entry point
- [ ] HTTPS access
- [ ] Cached proxy
- [ ] Run with flag configuration
- [ ] File based configuration

Apache License Version 2.0