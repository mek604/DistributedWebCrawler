package main

import (
	"net"
	"net/rpc"
	"log"
	"os"
	"strings"
	"sync"
)


var (
	workers []string // ip:port of workers
	WorkerDomainsList map[string][]string // maps each worker's ip:port to addresses
	mu sync.Mutex
)

//************************* RPC *************************
type LatencyStats struct {
	Min    int // min measured latency in milliseconds to host
	Median int // median measured latency in milliseconds to host
	Max    int // max measured latency in milliseconds to host
}

type MServer int
type MWorker int

func getPort (address string) (port string) {
	return	strings.Split(address, ":")[1]
}



//*************************************************************

/*

	known problems:


*/

func main() {

	args := os.Args[1:]
	
	if len(args) != 2 {
		log.Fatalf("server.go expected two arguments: [worker-incoming ip:port] [client-incoming ip:port]")
	}

	workerAddr := args[0]
	clientAddr := args[1]
	// workerPort := getPort(workerAddr)
	// clientPort := getPort(clientAddr)
	// assignedDomain := make(map[string][]string)

	// for Worker related RPC
	go func() {
		rpc.Register(new(MWorker))

		listener, err := net.Listen("tcp", workerAddr)
	    if err != nil {
	        log.Fatal("tcp server listener error", err)
	    }
	    for {
	        conn, err := listener.Accept()
	        if err != nil {
	            log.Fatal("tcp server accept error ", err)
	        }
	        log.Println("Worker Connected:", workerAddr)
			handleWorkerConnection(conn)
	    }
	}()

	// for Client related RPC
	go func() {
		rpc.Register(new(MServer))

		listener, err := net.Listen("tcp", clientAddr)
	    if err != nil {
	        log.Fatal("tcp server listener error ", err)
	    }
	    for {
	        conn, err := listener.Accept()
	        if err != nil {
	            log.Fatal("tcp server accept error ", err)
	        }
	        log.Println("Client Connected: ", clientAddr)
			handleClientConnection(conn)
	    }
	}()


	select {}
}


func handleClientConnection(conn net.Conn) {
	rpc.ServeConn(conn)
	conn.Close()

}


func handleWorkerConnection(conn net.Conn) {
	rpc.ServeConn(conn)
	workers = append(workers, conn.RemoteAddr().String())
	conn.Close()
}


//*************************************************************
type GetWorkersReq struct {
}

type GetWorkersRes struct {
	WorkerIPsList []string
}
func (t *MServer) GetWorkers(req *GetWorkersReq, reply *GetWorkersRes) error {
	reply.WorkerIPsList = workers
	return nil
}

//*************************************************************
type RegisterWorkerReq struct {
	IP string
}
type RegisterWorkerRes struct {}

func (t *MWorker) RegisterWorker(req *RegisterWorkerReq, reply *RegisterWorkerRes) error {
	mu.Lock()
	defer mu.Unlock()
	log.Println("Recieved RegisterWorker Request:", req.IP)
	workers = append(workers, req.IP)
	log.Println("List of workers:", workers)
	return nil
}
//*************************************************************

