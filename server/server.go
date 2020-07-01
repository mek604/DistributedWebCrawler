package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"github.com/tutorialedge/go-grpc-tutorial/chat"

)

const (
	port = ":9000"
)


func main() {
	lis, err := net.Listen("tcp",":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v " , port, err)
	}

	// generate in chat.pb.proto
	s := chat.Server{}

	grpcServer := grpc.NewServer()

	// register gRPC server , chat server struct
	chat.RegisterChatServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %v: %v", port, err)
	}

}