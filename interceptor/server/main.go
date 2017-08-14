package main

import (
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/chapin/go-rpc-tutorial/interceptor/proto"
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

// SayHello
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	resp := new(pb.HelloReply)
	resp.Messsage = "Hello " + in.Name + "."
	return resp, nil
}

// auth
func auth(ctx context.Context) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "无token认证信息")
	}
	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "10000" || appkey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "Token认证信息无效：appid=%s, appkey=%s", appid, appkey)
	}

	return nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("../keys/server.pem", "../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credential %v", err)
	}

	opts = append(opts, grpc.Creds(creds))

	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}
		return handler(ctx, req)
	}

	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...)

	pb.RegisterHelloServer(s, HelloService)

	grpclog.Println("Listen on " + Address + " with TLS + Token + Interceptor")

	s.Serve(listen)
}
