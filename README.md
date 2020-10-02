# protoc-gen-zap

**Warning: this is an early version - do not use in production**

Automated code generation for your protobuf objects to implement [`zap.ObjectMarshaler`](https://github.com/uber-go/zap/blob/master/zapcore/marshaler.go), based on [lyft/protoc-gen-star](https://github.com/lyft/protoc-gen-star)

This is useful to log protobuf objects in zap without heavy reflection:

```go
l, _ := zap.NewProduction()

l.Info("create-user",
  zap.Object("user", user)
)
```

## Requirements

- [protoc](http://google.github.io/proto-lens/installing-protoc.html)
- [go plugin for protoc](https://developers.google.com/protocol-buffers/docs/gotutorial)

## running tests

Code generation is done in the `protoc` flow:

```bash
go install . && protoc -I . -I ${GOPATH}/src --go_out=":./test" --zap_out="lang=go:./test" test/test.proto
```
