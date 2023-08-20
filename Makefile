.PHONY: generate

generate:
	protoc --go_out=./pkg/sandbox --go_opt=paths=source_relative \
    --go-grpc_out=./pkg/sandbox --go-grpc_opt=paths=source_relative \
    api/sandbox.proto