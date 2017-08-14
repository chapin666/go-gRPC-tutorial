package main

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc/credentials"

	pb "github.com/chapin/go-rpc-tutorial/interceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC server address
	Address = "127.0.0.1:50052"
	// OpenTLS 是否开启tls认证
	OpenTLS = true
)

type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "10000",
		"appkey": "i am key",
	}, nil
}
func (c customCredential) RequireTransportSecurity() bool {
	if OpenTLS {
		return true
	}
	return false
}

func main() {

	var err error
	var opts []grpc.DialOption

	if OpenTLS {
		creds, err := credentials.NewClientTLSFromFile("../keys/server.pem", "lycam")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %s", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	conn, err := grpc.Dial(Address, opts...)

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
