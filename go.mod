module github.com/opentofu/provider-client

go 1.26.0

require (
	github.com/google/go-cmp v0.7.0
	github.com/vmihailenco/msgpack/v5 v5.3.5
	github.com/zclconf/go-cty v1.16.3
	github.com/zclconf/go-cty-debug v0.0.0-20240509010212-0d6042c53940
	go.rpcplugin.org/rpcplugin v0.3.1
	go.uber.org/mock v0.6.0
	google.golang.org/genproto v0.0.0-20250715232539-7130f93afb79
	google.golang.org/grpc v1.81.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/apparentlymart/go-ctxenv v1.0.0 // indirect
	github.com/apparentlymart/go-shquot v0.0.1 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.6.2 // indirect
)

tool (
	go.uber.org/mock/mockgen
	golang.org/x/tools/cmd/stringer
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
