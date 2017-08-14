package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	pb "github.com/chapin/go-rpc-tutorial/hello-http/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

//
type helloHttpService struct{}

// HelloHttpService ...
var HelloHttpService = helloHttpService{}

func main() {
	endpoint := "127.0.0.1:50052"

	creds, err := credentials.NewClientTLSFromFile("../keys/server.pem", "../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	conn, _ := net.Listen("tcp", endpoint)
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterHelloHttpServer(grpcServer, HelloHttpService)

	ctx := context.Background()
	ctx, calcel := context.WithCancel(ctx)
	defer calcel()

	dcreds, err := credentials.NewClientTLSFromFile("../keys/server.pem", "lycam")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := runtime.NewServeMux()
	err = pb.RegisterHelloHttpHandlerFromEndpoint(ctx, gwmux, endpoint, dopts)
	if err != nil {
		fmt.Printf("Serve: %v\n", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/".gwmux)

	if err != nil {
		panic(err)
	}

	cert, _ := ioutil.ReadFile("../keys/server.pem")
	key, _ := ioutil.ReadFile("../keys/server.key")
	var demoKeyPair *tls.Certificate
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}
	demoKeyPair = &pair

	srv := &http.Server{
		Address: endpoint,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*demoKeyPair},
		}
	}

	fmt.Printf("grpc and https on port: %d\n", 50052)

	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}
