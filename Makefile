.PHONY: bin tools gen image push integration

OUTPUT_DIR=bin
OUTPUT_NAME=server
MAIN_FILE=cmd/server/main.go
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
PROTO_FILES=$(shell find . -name '*.proto' -type f -not -path "./vendor/*")

bin: gen $(GO_FILES)
	@ go build -o $(OUTPUT_DIR)/$(OUTPUT_NAME) $(MAIN_FILE)

tools:
	@ sudo apt update && sudo apt install -y protobuf-compiler
	@ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	@ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

gen: $(PROTO_FILES)
	@ mkdir -p gen
	@ protoc --go_out=gen --go_opt=paths=source_relative \
    --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
    api/protos/beacon.proto

image:
	@ docker build \
		-t troydai/grpcbeacon:latest \
		-t troydai/grpcbeacon:`git describe --tags` \
		.

push: image
	@ docker push troydai/grpcbeacon:`git describe --tags`

integration:
	@ ./scripts/integration-test.sh
