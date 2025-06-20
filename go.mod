module github.com/feight/newsteam-sdk

go 1.23.0

toolchain go1.23.3

require (
	buf.build/gen/go/dgroux/newsteam/connectrpc/go v1.18.1-20250620112026-bd6f53130430.1
	buf.build/gen/go/dgroux/newsteam/protocolbuffers/go v1.36.6-20250620112026-bd6f53130430.1
	connectrpc.com/connect v1.18.1
	github.com/fatih/color v1.18.0
	github.com/go-resty/resty/v2 v2.16.5
	github.com/google/uuid v1.6.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.41.0
)

require (
	buf.build/gen/go/dgroux/lib/protocolbuffers/go v1.36.6-20250620112026-22936ec2daf0.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
