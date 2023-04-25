FROM fullstorydev/grpcurl AS grpcurl
FROM golang:alpine3.17 AS builder

RUN apk update && apk add --no-cache make protobuf-dev protoc

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src

RUN mkdir -p gen
RUN make tools
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
COPY api/protos/beacon/api.proto /root/api.proto
COPY cmd/toolbox/* /root/

ENTRYPOINT ["tail", "-f", "/dev/null"]
