package main

import (
	"flag"
	"fmt"
	"github.com/huseyinbabal/grpc-ping-pong/pingpong"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var port = flag.Int("port", 10000, "Server port")

type pingPongServer struct {
	pingpong.UnimplementedPingPongServiceServer
}

func newPingPongServer() *pingPongServer {
	return &pingPongServer{}
}

func (s *pingPongServer) Ping(stream pingpong.PingPongService_PingServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&pingpong.PongResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()

	listen, listenErr := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if listenErr != nil {
		log.Fatalf("Failed to listen %v", listenErr)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pingpong.RegisterPingPongServiceServer(grpcServer, newPingPongServer())
	grpcServer.Serve(listen)
}
