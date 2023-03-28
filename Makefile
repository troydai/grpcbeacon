.PHONY: bin tools gen image push integration

# Override with setting these two and run make with option -e
ARCH=$(shell uname -m | tr '[:upper:]' '[:lower:]')
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')

OUTPUT_DIR=bin
OUTPUT_PATH=$(OUTPUT_DIR)/$(OS)/$(ARCH)
OUTPUT_NAME=server
MAIN_FILE=cmd/server/main.go
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
PROTO_FILES=$(shell find . -name '*.proto' -type f -not -path "./vendor/*")

bin: gen $(GO_FILES)
	@ echo "Building under $(OUTPUT_DIR) for $(OS)/$(ARCH)"
	@ GOOS=$(OS) GOARCH=$(ARCH) go build -v -o $(OUTPUT_PATH)/server   cmd/server/main.go
	@ GOOS=$(OS) GOARCH=$(ARCH) go build -v -o $(OUTPUT_PATH)/goclient cmd/goclient/main.go

run: bin
	@ $(OUTPUT_PATH)/server

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
	@ docker buildx build \
		--platform=arm64,amd64 \
		-t troydai/grpcbeacon:latest \
		-t troydai/grpcbeacon:`git describe --tags` \
		.

container-run: image
	@ docker run --platform=linux/arm64 --rm -it -p 50001:8080 troydai/grpcbeacon:latest

push: image
	@ docker push troydai/grpcbeacon:`git describe --tags`

integration:
	@ ./scripts/integration-test.sh
