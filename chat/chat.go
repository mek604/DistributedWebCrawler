// after generating .proto
// this is the RPC implementation

package chat

import (
	"log"
	"golang.org/x/net/context"
)

type Server struct {

}

// RPC method
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %s", message.Body)
	// return the msg to the client
	return &Message{Body: "Hello from the server"}, nil
}