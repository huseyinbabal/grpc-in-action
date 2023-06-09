# gRPC in Action
This is a sample demo project that I use within my live coding demo for presentation that you can see [here](https://bit.ly/grpc-in-action)

# Howto Run
1. `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pingpong/pingpong.proto`
2. `go run server/server.go`
3. `go run client/client.go`
