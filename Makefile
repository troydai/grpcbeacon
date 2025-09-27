.PHONY: bin tools gen image push integration

# Override with setting these two and run make with option -e
ARCH=$(shell uname -m | tr '[:upper:]' '[:lower:]')
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')

OUTPUT_DIR=bin
OUTPUT_NAME=server
MAIN_FILE=cmd/server/main.go
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
PROTO_FILES=$(shell find . -name '*.proto' -type f -not -path "./vendor/*")

tools:
	@ echo "Installing buf CLI..."
	@ go install github.com/bufbuild/buf/cmd/buf@v1.34.0

bin: gen $(GO_FILES)
	GOOS=$(OS) GOARCH=$(ARCH) go build -v -o $(OUTPUT_DIR)/$(OUTPUT_NAME) $(MAIN_FILE)

run: bin
	$(OUTPUT_DIR)/$(OUTPUT_NAME) -config=./demo/demo.conf

gen: $(PROTO_FILES)
	@ rm -rf gen/go
	@ buf generate

test:
	@ go test -v -race ./...