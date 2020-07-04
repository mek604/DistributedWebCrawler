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
<<<<<<< HEAD
	
=======
>>>>>>> 60fed5dee2a496342a4ec778900b33bc74ee8233
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

	log.Println("Sever accepting connections")

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
	log.Println("Registered to Worker Pool:", req.WorkerIP)
<<<<<<< HEAD
=======
	return nil
}
//*************************************************************
type DomainsReq struct {
	WorkerIP string // IP of worker
}
type DomainsRes struct {
	Domains []string // List of domain string
}
func (m *MServer) Domains(req *DomainsReq, res *DomainsRes) error {
	res.Domains = assignedDomains[req.WorkerIP]
	return nil
}
//*************************************************************
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
func hasMappedDomain(uri string) string {
	for ip, domains := range assignedDomains {
		for _,d := range domains {
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
func (m *MServer) Crawl(req *CrawlReq, res *CrawlRes) error {
	log.Println("Start crawling", req.URL)
	latReq := MeasureLatencyReq{
		URL: req.URL,
		Samples: samples,
	}
	latRes := MeasureLatencyRes{}
	stats := make(map[string]MeasureLatencyRes)

	worker := hasMappedDomain(req.URL)
	if worker == "" {
		var wg sync.WaitGroup
		wg.Add(len(workers))
		for _, worker := range workers {
			go func() {
				conn, err := net.Dial("tcp", worker)
				defer conn.Close()
				if err != nil {
					log.Fatal("Crawl: tcp dial error", err)
				}
				client := rpc.NewClient(conn)
				err = client.Call("MWorker.MeasureLatency", &latReq, &latRes)
				// handle error
				// log.Println(latReq.URL, latRes.Min, latRes.Median, latRes.Max)
				mu.Lock()
				stats[worker] = latRes
				mu.Unlock()
			}()
		}
		wg.Wait()
		log.Println("Mapped")
		worker := selectBestLatency(stats)
	} else {
		// add to domain
		assignedDomains[worker] = append(assignedDomains[worker], req.URL)
	}
	log.Println("Assigend", worker, "to", req.URL)
	
	conn, _ := net.Dial("tcp", worker)
	client := rpc.NewClient(conn)
	_ = client.Call("MWorker.CrawlWebsite", &latReq, &latRes)
	conn.Close()



>>>>>>> 60fed5dee2a496342a4ec778900b33bc74ee8233
	return nil
}
//*************************************************************
type DomainsReq struct {
	WorkerIP string // IP of worker
}
type DomainsRes struct {
	Domains []string // List of domain string
}
func (m *MServer) Domains(req *DomainsReq, res *DomainsRes) error {
	res.Domains = assignedDomains[req.WorkerIP]
	return nil
}
//*************************************************************
type CrawlReq struct {
	URL   string // URL of the website to crawl
	Depth int    // Depth to crawl to from URL
}
// Response to MServer.Crawl
type CrawlRes struct {
	WorkerIP string // workerIP
}

<<<<<<< HEAD
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
func hasMappedDomain(uri string) string {
	for ip, domains := range assignedDomains {
		for _,d := range domains {
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
func (m *MServer) Crawl(req *CrawlReq, res *CrawlRes) error {
	log.Println("Start crawling", req.URL)
	latReq := MeasureLatencyReq{
		URL: req.URL,
		Samples: samples,
	}
	latRes := MeasureLatencyRes{}
	stats := make(map[string]MeasureLatencyRes)

	worker := hasMappedDomain(req.URL)
	if worker == "" {
		var wg sync.WaitGroup
		wg.Add(len(workers))
		// request workers to measure latency to a website
		for _, worker := range workers {
			go func() {
				conn, err := net.Dial("tcp", worker)
				defer conn.Close()
				if err != nil {
					log.Fatal("Crawl: tcp dial error", err)
				}
				client := rpc.NewClient(conn)
				err = client.Call("MWorker.MeasureLatency", &latReq, &latRes)
				// handle error
				// log.Println(latReq.URL, latRes.Min, latRes.Median, latRes.Max)
				mu.Lock()
				stats[worker] = latRes
				mu.Unlock()
			}()
		}
		wg.Wait()
		log.Println("Mapped")
		worker := selectBestLatency(stats)
	} else {
		assignedDomains[worker] = append(assignedDomains[worker], req.URL)
	}
	log.Println("Assigend", worker, "to", req.URL)
	sendCrawlTask(worker)	
		
	return nil
}
//*************************************************************

func sendCrawlTask(worker string, ) {
	conn, _ := net.Dial("tcp", worker)
	client := rpc.NewClient(conn)

	// need request and response
	_ = client.Call("MWorker.CrawlWebsite", &latReq, &latRes)
	conn.Close()
}


=======
>>>>>>> 60fed5dee2a496342a4ec778900b33bc74ee8233
//*************************************************************
