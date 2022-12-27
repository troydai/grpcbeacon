#!/bin/bash

docker build . -t troydai/grpcecho:integration-test

docker run -d -P -l "purpose=test" troydai/grpcecho:integration-test > /dev/null

PORT=`docker ps -f label=purpose=test -qa | xargs docker inspect | jq -r '.[0].NetworkSettings.Ports."8080/tcp" | .[0].HostPort'`

echo "Port is $PORT"

grpcurl -plaintext "localhost:$PORT" describe

grpcurl -plaintext -d '{"greet": "sample"}' "localhost:$PORT" Service.Echo 

docker ps -f label=purpose=test -qa | xargs docker rm -f > /dev/null
