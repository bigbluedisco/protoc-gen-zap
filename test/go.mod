module github.com/bigbluedisco/protoc-gen-zap/test

go 1.26.0

require (
	github.com/bigbluedisco/protoc-gen-zap v0.0.0
	go.uber.org/zap v1.15.0
	google.golang.org/genproto v0.0.0-20261005182115-fad411399dd8
	google.golang.org/protobuf v1.36.12
)

require (
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/bigbluedisco/protoc-gen-zap => ../
