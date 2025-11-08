.PHONY: proto-go

# Build .proto for Go
proto-go:
	@mkdir -p pkg/grpc
	@protoc -I api/proto --go_out=pkg/grpc --go_opt=paths=source_relative --go-grpc_out=pkg/grpc --go-grpc_opt=paths=source_relative api/proto/api.proto
