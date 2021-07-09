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

To obfsucate a field, you may use an annotation like in the following example
```proto
message Test {
    string id = 1 [(zap.obfuscation_type) = STARS]; // Replace the value with 3 stars.
    string secret_id = 2 [(zap.obfuscation_type) = HIDE]; // Removes the field from the logs.
    map<string, string> my_hidden_map = 3 [(zap.obfuscation_type) = STARS];
    repeated uint32 uint_array = 4 [(zap.obfuscation_type) = STARS];
}
```

The result 
```go
func (c Test) MarshalLogObject(o zapcore.ObjectEncoder) error {

	o.AddString("id", "***")

	for key := range c.GetMyHiddenMap() {
		  o.AddString("my_hidden_map_"+key, "***")
	}

	uintarrayLength := len(c.GetCoucou())
	o.AddInt("coucou_length", uintarrayLength)

	if uintarrayLength > 100 {
		uintarrayLength = 100
	}

	uintarray := make(utils.StringArray, uintarrayLength)
	for i := 0; i < uintarrayLength; i++ {
		uintarray[i] = "***"
		continue
	}

	return nil
}
```