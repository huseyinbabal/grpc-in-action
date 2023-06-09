package main

import (
	"context"
	"flag"
	"github.com/huseyinbabal/grpc-ping-pong/pingpong"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "gRPC server address")
	messages   = []string{
		"Hello!",
		"Welcome!",
		"Goodbye!",
		"Greetings!",
		"Have a nice day!",
	}
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("failed to dial %v", err)
	}

	defer conn.Close()

	client := pingpong.NewPingPongServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	stream, streamErr := client.Ping(ctx)

	if streamErr != nil {
		log.Fatalf("Failed to stream: %v", streamErr)
	}

	done := make(chan bool)
	streamContext := stream.Context()
	go func() {

		for {
			msg := generateRandomMessage()
			req := pingpong.PingRequest{Message: msg}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("Cannot send %v", req)
			}
			time.Sleep(time.Second * 2)
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
			log.Println(response.Message)
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

func generateRandomMessage() string {
	index := rand.Intn(len(messages)) // Generate random index
	return messages[index]            // Return random message
}
