package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"log"
	"strings"
)

// server RPC
type MWorker int
// worker RPC
type WWorker int


var serverAddress string
var workerPort string

type RegisterWorkerRes struct {}
type RegisterWorkerReq struct {
	IP string
}


func getAddress (address string) (string) {
	return	strings.Split(address, ":")[0]
}


func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Usage: go run worker.go [server ip:port] [worker port]")
		return
	}

	serverAddr := args[0]
	// workerPort := args[1]

	go func() {
		workerServer := rpc.NewServer()
		workerServer.Register(new(MWorker))

		listener, err := net.Listen("tcp", ":" + workerPort)
		if err != nil {
			log.Fatal("tcp server listener error ", err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("tcp server accept error ", err)
			}
			log.Println("Connected to server", serverAddr)
			go workerServer.ServeConn(conn)
		}
	}()

	select{}
}

func handleClientConnection(conn net.Conn, serverAddr string) {
	rpc.ServeConn(conn)
}