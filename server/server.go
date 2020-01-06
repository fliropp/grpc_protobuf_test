package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	ping "github.com/fliropp/grpc_protobuf_test/ping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type pingServer struct {
	msg string
}

func (ps *pingServer) GetPing(ctx context.Context, p *ping.PingReq) (*ping.PingResp, error) {
	return &ping.PingResp{Response: "pong"}, nil
}

func (s *pingServer) ListFeatures(p *ping.PingReq, stream ping.Ping_StreamPingServer) error {
	pings := [5]string{"ping1", "ping2", "ping3", "ping4", "ping5"}
	for _, p := range pings {
		err := stream.Send(&ping.PingResp{Response: p})
		if err != nil {
			return err
		}
	}
	return nil
}

func newServer() *pingServer {
	s := &pingServer{msg: "ping"}
	return s
}

func main() {
	fmt.Println("gRPC ping server is running")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	ping.RegisterPingServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
