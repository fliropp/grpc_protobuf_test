package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	ping "github.com/fliropp/grpc_protobuf_test/ping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func streamPings(client ping.PingClient, p *ping.PingReq) {
	log.Printf("Sedning msg ovre gRPC (%s)", p)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.StreamPing(ctx, p)
	if err != nil {
		log.Fatalf("%v.GetPingReq(_) = _, %v: ", client, err)
	}
	for {
		png, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.StreamPings(_) = _, %v", client, err)
		}
		log.Println(png)
	}
}

func getSinglePing(client ping.PingClient, p *ping.PingReq) {
	log.Printf("Getting single ping over gRPC (%s)", p)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	png, err := client.GetPing(ctx, p)
	if err != nil {
		log.Fatalf("%v.GetPingReq(_) = _, %v: ", client, err)
	}
	log.Println(png)
}

func main() {
	fmt.Println("gRPC ping client is running")
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := ping.NewPingClient(conn)

	getSinglePing(client, &ping.PingReq{Request: "ping"})
	fmt.Println("-----------")
	streamPings(client, &ping.PingReq{Request: "pings"})

}
