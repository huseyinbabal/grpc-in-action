package main

import (
	"context"
	"flag"
	"github.com/huseyinbabal/grpc-in-action/repository"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

var serverAddr = flag.String("server_addr","localhost:10000", "gRPC server address")
func main() {
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts,grpc.WithInsecure())
	opts=append(opts,grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err!=nil{
		log.Fatalf("failed to dial %v",err)
	}

	defer conn.Close()

	client := repository.NewRepositoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, streamErr := client.Search(ctx)

	if streamErr != nil {
		log.Fatalf("Failed to stream: %v", streamErr)
	}

	keywords := []string{"raft", "spring"}
	done := make(chan bool)
	streamContext := stream.Context()
	go func() {
		for _, keyword := range keywords {
			req := repository.SearchCodeRequest{User: "huseyinbabal", Keyword: keyword}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("Cannot sned %v", req)
			}
			time.Sleep(time.Second * 2)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("Couldn't close send: %v", err)
		}
	}()

	go func() {
		for {
			response, respErr := stream.Recv()
			if respErr == io.EOF {
				close(done)
				return
			}
			if respErr != nil {
				log.Fatalf("Couldn't receive %v", respErr)
			}
			log.Println(response.Name)
		}
	}()

	go func() {
		<-streamContext.Done()
		if err := streamContext.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()
	<-done
	log.Println("Streaming finished")
}
