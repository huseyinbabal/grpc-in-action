package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v38/github"
	"github.com/huseyinbabal/grpc-in-action/repository"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var port = flag.Int("port", 10000, "Server port")

type repositoryServer struct {
	repository.UnimplementedRepositoryServer
	githubClient *github.Client
}

func newRepositoryServer() *repositoryServer {
	return &repositoryServer{githubClient: github.NewClient(nil)}
}

func (s *repositoryServer) Search(stream repository.Repository_SearchServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF{
			return nil
		}

		if err != nil{
			return err
		}

		page := 1
		lastPage := false

		for !lastPage{
			opts:=&github.SearchOptions{
				Sort:        "forks",
				Order:       "desc",
				ListOptions: github.ListOptions{
					Page: page,
					PerPage: 10,
				},
			}
			searchResponse, _, searchErr := s.githubClient.Search.Code(context.Background(), "q="+in.Keyword+" user:"+in.User, opts)
			if searchErr!=nil{
				return searchErr
			}

			if len(searchResponse.CodeResults)==0{
				lastPage=true
			}
			page+=1
			for _, codeResult := range searchResponse.CodeResults {
				if err :=stream.Send(&repository.SearchCodeResponse{Name: *codeResult.Name});err !=nil{
					return err
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	listen, listenErr := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if listenErr != nil{
		log.Fatalf("Failed to listen %v",listenErr)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	repository.RegisterRepositoryServer(grpcServer, newRepositoryServer())
	grpcServer.Serve(listen)
}
