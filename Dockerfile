FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx

FROM --platform=$BUILDPLATFORM golang:alpine AS builder

COPY --from=xx / /
ARG TARGETPLATFORM

RUN xx-info env

RUN xx-apk update && xx-apk add --no-cache make protobuf-dev
RUN xx-go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN xx-go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN xx-go mod download

RUN mkdir -p gen
RUN protoc --go_out=gen --go_opt=paths=source_relative \
    --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
    api/protos/beacon.proto

COPY . /src
RUN xx-go build -v -o bin/server   cmd/server/main.go
RUN xx-go build -v -o bin/goclient cmd/goclient/main.go

FROM --platform=$BUILDPLATFORM scratch AS server

WORKDIR /opt/bin
COPY --from=builder /src/bin/server /opt/bin/server

EXPOSE 8080

ENTRYPOINT [ "/opt/bin/server" ]

FROM --platform=$BUILDPLATFORM scratch AS prober

WORKDIR /opt/bin
COPY --from=builder /src/bin/goclient /opt/bin/goclient

EXPOSE 8080

ENTRYPOINT [ "/opt/bin/goclient" ]
