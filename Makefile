.PHONY: bin tools gen image push integration setup check-tools

# Override with setting these two and run make with option -e
ARCH=$(shell uname -m | tr '[:upper:]' '[:lower:]' | sed 's/x86_64/amd64/')
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')

OUTPUT_DIR=bin
OUTPUT_NAME=server
MAIN_FILE=cmd/server/main.go
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
PROTO_FILES=$(shell find . -name '*.proto' -type f -not -path "./vendor/*")

bin: gen $(GO_FILES)
	GOOS=$(OS) GOARCH=$(ARCH) go build -v -o $(OUTPUT_DIR)/$(OUTPUT_NAME) $(MAIN_FILE)

run: bin
	$(OUTPUT_DIR)/$(OUTPUT_NAME) -config=./demo/demo.conf

check-tools:
	@which buf >/dev/null || (echo "❌ buf not found. Run 'make setup' or install with: go install github.com/bufbuild/buf/cmd/buf@latest" && exit 1)
	@echo "✅ All required tools are available"

setup:
	@echo "🔧 Running development environment setup..."
	@./scripts/setup-dev.sh

gen: check-tools $(PROTO_FILES)
	@ rm -rf gen/go
	@ buf generate
	@echo "✅ Protocol buffer code generated successfully"

test:
	@ go test -v -race ./...