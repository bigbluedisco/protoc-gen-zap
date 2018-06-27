# protoc-gen-zap

Automated code generation for your protobuf objects to implement `zap.ObjectMarshaler`, using [lyft/protoc-gen-star](https://github.com/lyft/protoc-gen-star)

This is useful to log protobuf objects in zap without heavy reflection:

``` go
l, _ := zap.NewProduction()

l.Info("create-user",
  zap.Object("user", user)
)
```

## installation

Make sure `dep` is installed on your machine, then run:

``` bash
dep ensure
go build
go install
```

## running

Code generation is done in the `protoc` flow:

``` bash
protoc -I . -I ${GOPATH}/src --go_out=":./out" --zap_out="lang=go:./out" test.proto
```
