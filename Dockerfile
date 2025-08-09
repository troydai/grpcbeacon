FROM fullstorydev/grpcurl AS grpcurl
FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache make protobuf-dev protoc

# Install buf CLI for code generation
RUN go install github.com/bufbuild/buf/cmd/buf@v1.31.0

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src

RUN mkdir -p gen
# Removed: RUN make tools (target does not exist)
RUN make gen

RUN go build -v -o bin/server cmd/server/main.go

FROM scratch AS server

WORKDIR /opt/bin
COPY --from=builder /src/bin/server /opt/bin/server

EXPOSE 8080

ENTRYPOINT [ "/opt/bin/server" ]

FROM alpine AS toolbox

RUN apk add curl bash jq vim tcpdump

COPY --from=grpcurl /bin/grpcurl /bin/grpcurl
# Copy the correct proto file path
COPY proto/troydai/grpcbeacon/v1/api.proto /root/api.proto
COPY cmd/toolbox/* /root/

ENTRYPOINT ["tail", "-f", "/dev/null"]
