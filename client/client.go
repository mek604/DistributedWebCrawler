package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/tutorialedge/go-grpc-tutorial/chat"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	defer conn.Close()
	c := chat.NewChatServiceClient(conn)

	// can interact with gRPC server 
	msg := chat.Message {
		Body: "Hello from the client!",
	}
	response, err := c.SayHello(context.Background(), &msg)
	if nil != err {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)
}