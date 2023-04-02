#!/bin/bash

set -e

IMAGE_NAME="troydai/grpcbeacon"
IMAGE_TAG="integration-test"
IMAGE="$IMAGE_NAME:$IMAGE_TAG"
ARCH=`uname -m | tr '[:upper:]' '[:lower:]'`
OS=`uname -s | tr '[:upper:]' '[:lower:]'`

echo "[1] Building"
docker build . -q -t $IMAGE > /dev/null
make bin

echo "[2] Starting docker container"
docker ps -f label=purpose=test -qa | xargs docker rm -f > /dev/null
docker run -d -P -l "purpose=test" $IMAGE > /dev/null
PORT=`docker ps -f label=purpose=test -qa | xargs docker inspect | jq -r '.[0].NetworkSettings.Ports."8080/tcp" | .[0].HostPort'`

echo "    Container label purpose=test"
echo "    Container port $PORT"

echo "[3] Testing describe gRPC interface"
grpcurl -plaintext "localhost:$PORT" describe > /dev/null

echo "[4] Testing beacon signal"
grpcurl -plaintext -d '' "localhost:$PORT" grpcbeacon.Beacon/Signal > /dev/null

echo "[5] Testing Go gRPC client"
bin/$OS/$ARCH/goclient "localhost:$PORT" > /dev/null

echo "[0] Clean up containers"
docker ps -f label=purpose=test -qa | xargs docker rm -f > /dev/null
