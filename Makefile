.PHONE: tools gen

bin:
	@ go build -o artifacts/server cmd/server/main.go 

tools:
	@ sudo apt update && sudo apt install -y protobuf-compiler
	@ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	@ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

gen:
	@ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protos/echo.proto

image:
	@ docker build . -t troydai/grpcecho:latest