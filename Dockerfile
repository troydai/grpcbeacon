FROM fullstorydev/grpcurl AS grpcurl
FROM golang:alpine3.17 AS builder

RUN apk update && apk add --no-cache make protobuf-dev protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src

RUN mkdir -p gen
RUN protoc --go_out=gen --go_opt=paths=source_relative \
    --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
    api/protos/beacon.proto

RUN go build -v -o bin/server   cmd/server/main.go
RUN go build -v -o bin/goclient cmd/goclient/main.go

FROM scratch AS server

WORKDIR /opt/bin
COPY --from=builder /src/bin/server /opt/bin/server

EXPOSE 8080

ENTRYPOINT [ "/opt/bin/server" ]

FROM scratch AS prober

WORKDIR /opt/bin
COPY --from=builder /src/bin/goclient /opt/bin/goclient

EXPOSE 8080

ENTRYPOINT [ "/opt/bin/goclient" ]

FROM alpine AS toolbox

RUN apk add curl bash jq

COPY --from=grpcurl /bin/grpcurl /bin/grpcurl
COPY --from=builder /src/api/protos/beacon.proto /etc/protos/beacon.proto
COPY --from=builder /src/cmd/toolbox/curl.sh /run/curl.sh

ENTRYPOINT ["tail", "-f", "/dev/null"]
