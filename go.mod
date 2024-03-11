module github.com/feight/newsteam-sdk

go 1.22.0

require (
	buf.build/gen/go/dgroux/newsteam/connectrpc/go v1.15.0-20240308130319-b9d3bc9b0677.1
	buf.build/gen/go/dgroux/newsteam/protocolbuffers/go v1.33.0-20240308130319-b9d3bc9b0677.1
	connectrpc.com/connect v1.15.0
	github.com/google/uuid v1.6.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.22.0
	google.golang.org/protobuf v1.33.0
)

require (
	buf.build/gen/go/envoyproxy/protoc-gen-validate/protocolbuffers/go v1.33.0-20231130202533-71881f09a0c5.1 // indirect
	golang.org/x/text v0.14.0 // indirect
)
