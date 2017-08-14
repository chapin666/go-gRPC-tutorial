#! /bin/bash

# 编译google api，新版编译器可以省略M参数
protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto

# 编译hello_http.proto
protoc -I . --go_out=plugins=grpc,Mgoogle/api/annotations.proto=github.com/chapin/go-rpc-tutorial/hello-http/proto/google/api:. *.proto

# 编译hello_http.proto gateway
protoc --grpc-gateway_out=logtostderr=true:. hello_http.proto

# test
curl -X POST -k https://localhost:50052/example/echo -d '{"name": "gRPC-HTTP is working!"}'