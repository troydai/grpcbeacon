#!/bin/bash

IMAGE_NAME="troydai/grpcbeacon"
IMAGE_TAG="integration-test"
IMAGE="$IMAGE_NAME:$IMAGE_TAG"

docker build . -t $IMAGE

docker run -d -P -l "purpose=test" $IMAGE > /dev/null

PORT=`docker ps -f label=purpose=test -qa | xargs docker inspect | jq -r '.[0].NetworkSettings.Ports."8080/tcp" | .[0].HostPort'`

echo "Port is $PORT"

grpcurl -plaintext "localhost:$PORT" describe

grpcurl -plaintext -d '' "localhost:$PORT" grpcbeacon.Beacon/Signal 

docker ps -f label=purpose=test -qa | xargs docker rm -f > /dev/null
