# goorita
Simple http proxy written in Go

## limitation
- Binding to Network Interface only work in Unix/Linux env since it needs `SO_BINDTODEVICE` 

## features
- [x] Multi outgoing IP address :checked
- [x] Multi port entry point :checked
- [ ] HTTPS access :unchecked
- [ ] Cached proxy :unchecked

Apache License Version 2.0