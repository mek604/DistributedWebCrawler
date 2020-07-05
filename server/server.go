package main

import (
	"net"
	"net/rpc"
	"net/url"
	"log"
	"os"
	"strings"
	"sync"
	"strconv"
	"time"
)


var (
	workers []string // ip:port of workers
	assignedDomains map[string][]string // maps each worker's ip:port to addresses
	samples int // number of samples used to generate RTT stats
	mu sync.Mutex
)

//************************* RPC *************************
type MServer int
type MWorker int

func getPort (address string) (port string) {
	return	strings.Split(address, ":")[1]
}



//*************************************************************

/*



*/

func main() {

	args := os.Args[1:]

	if len(args) != 3 {
		log.Fatalf("server.go expected two arguments: [worker-incoming ip:port] [client-incoming ip:port] [# of RTT sample]")
	}

	workerAddr := args[0]
	clientAddr := args[1]
	samples,_ = strconv.Atoi(args[2])
	assignedDomains = make(map[string][]string)

	log.Println("Accepting connections")

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
	        log.Println("Client Connected:", clientAddr)
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
	conn.Close()
}
//*************************************************************
type GetWorkersReq struct {}
type GetWorkersRes struct {
	WorkerIPsList []string
}
// Retuns a list of worker addresses [ip:port]
func (t *MServer) GetWorkers(req *GetWorkersReq, res *GetWorkersRes) error {
	res.WorkerIPsList = workers
	return nil
}

//*************************************************************
type RegisterWorkerReq struct {
	WorkerIP string
}
type RegisterWorkerRes struct {}
func (t *MWorker) RegisterWorker(req *RegisterWorkerReq, res *RegisterWorkerRes) error {
	mu.Lock()
	defer mu.Unlock()
	workers = append(workers, req.WorkerIP)
	log.Println("Register", req.WorkerIP, " to worker pool")
	return nil
}
//*************************************************************
//*************************************************************
type DomainsReq struct {
	WorkerIP string // IP of worker
}
type DomainsRes struct {
	Domains []string // List of domain string
}
func (m *MServer) Domains(req *DomainsReq, res *DomainsRes) error {
	log.Println("<- Client RPC: returning a map of a worker to domains")
	res.Domains = assignedDomains[req.WorkerIP]
	return nil
}
//*************************************************************
// client's crawl request and response
type CrawlReq struct {
	URL   string // URL of the website to crawl
	Depth int    // Depth to crawl to from URL
}
// Response to MServer.Crawl
type CrawlRes struct {
	WorkerIP string // workerIP
}
// measuring latency between a website and a worker
type MeasureLatencyRes struct {
	Min time.Duration
	Median time.Duration
	Max time.Duration
}
type MeasureLatencyReq struct {
	URL string
	Samples int
}
func getDomainName(uri string) string {
	u, _ := url.Parse(uri)
	return u.Host
}
func findMappedDomain(uri string) string {
	for ip, domains := range assignedDomains {
		for _,d := range domains {
			// log.Println("\t\t", uri, getDomainName(uri), d)
			if getDomainName(uri) == d {
				return ip
			}
		}
	}
	return ""
}
func selectBestLatency(stats map[string]MeasureLatencyRes) (currIP string) {
	minLatency,_ := time.ParseDuration("24h")
	currIP = ""
	for ip, latency := range stats {
		if latency.Min < minLatency {
			minLatency = latency.Min
			currIP = ip
		}
	}
	return
}
// find a worker that owns the domain of a URL:
// if owner not found, ask all workers to measure latency to the URL
// and select the worker with lowest latency
func findWorker(reqURL string) string {
	latReq := MeasureLatencyReq{
		URL: reqURL,
		Samples: samples,
	}
	latRes := MeasureLatencyRes{}
	stats := make(map[string]MeasureLatencyRes)
	worker := findMappedDomain(reqURL)
	if worker == "" {
		log.Printf("Owner-worker of (%s) domain not found\n", reqURL)
		var wg sync.WaitGroup
		wg.Add(len(workers))
		// request workers to measure latency to a website
		log.Printf("\tMeasuring latency of (%s)", reqURL)
		for _, worker := range workers {
			go func() {
				conn, err := net.Dial("tcp", worker)
				defer conn.Close()
				if err != nil {
					log.Fatal("Crawl: tcp dial error", err)
				}
				client := rpc.NewClient(conn)
				err = client.Call("MWorker.MeasureLatency", latReq, &latRes)
				// handle error
				// log.Println(latReq.URL, latRes.Min, latRes.Median, latRes.Max)
				stats[worker] = latRes
				wg.Done()
			}()
		}
		wg.Wait()
		worker = selectBestLatency(stats)
		log.Printf("Assigned %s the domain of %s\n", worker, reqURL)
		mu.Lock()
		assignedDomains[worker] = append(assignedDomains[worker], getDomainName(reqURL))
		mu.Unlock()
	} 
	return worker
}
func (m *MServer) Crawl(req *CrawlReq, res *CrawlRes) error {
	log.Printf("<- Clien RPC: Crawl (%s)\n", req.URL)
	worker := findWorker(req.URL)
	res.WorkerIP = worker

	if req.Depth > 0 {
		sendCrawlTask(worker, req.URL, req.Depth)	
	}

	return nil
}
//*************************************************************
// worker's crawl request and response
type CrawlWebsiteReq struct {
	URL string
	Depth int
}
type CrawlWebsiteRes struct {
	Depth int
	Links []string
}
// does not seem to loop back
// maybe the crawl loop returns first?
func sendCrawlTask(worker, uri string, depth int) {
	log.Printf("-> Worker RPC: ask %s to crawl %s\n", worker, uri)
	conn, _ := net.Dial("tcp", worker)
	client := rpc.NewClient(conn)
	defer conn.Close()

	req := CrawlWebsiteReq{
		URL: uri,
		Depth: depth,
	}
	res := CrawlWebsiteRes{}
	// need to make this better
	_ = client.Call("MWorker.CrawlWebsite", req, &res)
	depth--
	var wg sync.WaitGroup
	// subsequent crawls
	if depth > 0 {
		for _, link := range res.Links {
			wg.Add(1)
			go func(link string, d int) {
				defer wg.Done()
				if link == "" {
					return
				}
				worker := findWorker(link)
				log.Println(link, worker, depth)
				sendCrawlTask(worker, link, d)
			}(link, depth)
		}
		depth--
		wg.Wait()
	}
}
//*************************************************************
