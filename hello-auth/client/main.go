package main

import (
	"context"

	"google.golang.org/grpc/credentials"

	pb "github.com/chapin/go-rpc-tutorial/hello-auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC server address
	Address = "127.0.0.1:50052"
)

func main() {

	creds, err := credentials.NewClientTLSFromFile("../keys/server.pem", "lycam")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %s", err)
	}

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))

	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()
	c := pb.NewHelloClient(conn)

	reqBody := new(pb.HelloRequest)
	reqBody.Name = "gRPC"
	r, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}
	grpclog.Println(r.Messsage)
}
