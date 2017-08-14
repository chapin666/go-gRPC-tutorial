package main

import (
	"net"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/chapin/go-rpc-tutorial/hello-auth/proto"
	"golang.org/x/net/context"
)

const (
	// Address gRPC server address
	Address = "127.0.0.1:50052"
)

// 定义helloSerice并实现约定接口

type helloService struct{}

// HelloService instance
var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	resp := new(pb.HelloReply)
	resp.Messsage = "Hello " + in.Name + "."

	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("../keys/server.pem", "../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credential %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterHelloServer(s, HelloService)

	grpclog.Println("Listen on " + Address + " with TLS")

	s.Serve(listen)
}
