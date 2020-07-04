package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/url"
	"net/http"
	"os"
	"log"
	"strings"
	"time"
	"sort"
)

// server RPC
type WServer int
// worker RPC
type MWorker int


const workerPort string = "3800"
var serverAddress string

type RegisterWorkerReq struct {
	WorkerIP string
}
type RegisterWorkerRes struct {}


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
			go workerServer.ServeConn(conn)
		}
	}()

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("tcp server dial error ", err)
	}
	client := rpc.NewClient(conn)
	registerAddr := getAddress(serverAddr) + ":" + workerPort
	req := RegisterWorkerReq{WorkerIP: registerAddr,}
	res := RegisterWorkerReq{}
	err = client.Call("MWorker.RegisterWorker", &req, &res)
	if err != nil {
		log.Fatal("rcp call error ", err)
	}
	conn.Close()

	select{}
}

func getDomainName(uri string) string {
	u, _ := url.Parse(uri)
	return u.Host
}
func getAbsolutePath(uri string) string {
	u, _ := url.Parse(uri)
	base, _ := url.Parse(getDomainName("http://example.com/directory/"))
	return base.ResolveReference(u).String()
}

func handleClientConnection(conn net.Conn, serverAddr string) {
	rpc.ServeConn(conn)
}

//*************************************************************
type MeasureLatencyReq struct {
	URL string
	Samples int
}
type MeasureLatencyRes struct {
	Min time.Duration
	Median time.Duration
	Max time.Duration
}
func (t *MWorker) MeasureLatency(req *MeasureLatencyReq, res *MeasureLatencyRes) error {
	var latencies []time.Duration
	skips := 0

	for i:=0 ; i<req.Samples; i++ {
		start := time.Now()
		response, err := http.Get(req.URL) 
		defer response.Body.Close()
		if err != nil {
			skips++
			log.Println("Http Get Failed ", err)
			continue
		}
		end := time.Now()
		elapsed := end.Sub(start)
		if response.StatusCode == 200 {
			latencies = append(latencies, elapsed)
		} else {
			skips++
		}
		fmt.Println("HTTP Request", req.URL, elapsed)
	}

	// sort latencies
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] > latencies[j]
	})
	// log.Println(latencies)
	if len(latencies) > 0 {
		res.Max = latencies[len(latencies) - 1]
		res.Min = latencies[0]
		n := len(latencies) - skips
		if n % 2 == 0 {
			res.Median = latencies[(n-1)/2] + latencies[(n-1/2) + 1] / 2
		} else {
			res.Median = latencies[(n - 1) / 2]
		}
	}
	return nil
}

//*************************************************************
type CrawlReq struct {
}
type CrawlRes struct {
}
func (t *MWorker) CrawlWebsite(req *CrawlReq, res *CrawlRes) error {

	return nil
}

//*************************************************************
