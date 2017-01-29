#!/bin/bash

set -e

if [ -z "$1" ]; then
	echo "Error: you must specify a .proto file to generate proxy for" >/dev/stderr
	exit 1

fi

PROTOFILE=$1

# generate go grpc files
protoc -I/usr/local/lib -I/usr/local/include -I.  -I$GOPATH/src  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:.  $PROTOFILE 


# generate rest to grpc gateway
protoc -I/usr/local/include -I.  -I$GOPATH/src  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --grpc-gateway_out=logtostderr=true:. $PROTOFILE 


# generate swagger file

protoc -I/usr/local/include -I.  -I$GOPATH/src  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --swagger_out=logtostderr=true:. $PROTOFILE 



