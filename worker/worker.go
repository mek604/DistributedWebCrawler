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
	"golang.org/x/net/html"
	// "sync"
)

// server RPC
type WServer int
// worker RPC
type MWorker int


const workerPort string = "3800"
var (
	serverAddress string
	savedDomains map[string]struct{}
	webGraph map[string][]string // maps url:[next urls]
	// mu sync.Mutex
)

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

	serverAddress = args[0]
	webGraph = make(map[string][]string)
	savedDomains = make(map[string]struct{})

	log.Println("Connected to server", serverAddress)

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

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal("tcp server dial error ", err)
	}
	client := rpc.NewClient(conn)
	registerAddr := getAddress(serverAddress) + ":" + workerPort
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
// remove relative path from a uri starting with "http"
func resolveReference(uri string) string {
	u, _ := url.Parse(uri)
	base, _ := url.Parse(getDomainName(uri))
	return base.ResolveReference(u).String()
}
func filterAddress(link, domain string) string {
	resolved := resolveReference(link)
	if strings.HasPrefix(link, "/") {
		return domain + resolved
	} 
	return resolved
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
			log.Println("HTTP Get Failed ", err)
			continue
		}
		end := time.Now()
		elapsed := end.Sub(start)
		if response.StatusCode == 200 {
			latencies = append(latencies, elapsed)
		} else {
			skips++
		}
		log.Printf("Measuring HTTP Get (%s) = %v\n", req.URL, elapsed)
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
type CrawlWebsiteReq struct {
	URL string
	Depth int
}
type CrawlWebsiteRes struct {
	Depth int
	Links []string
}
func contains(arr []string, val string) bool {
	for _, a := range arr {
		if val == a {
			return true
		}
	}
	return false
}
func (t *MWorker) CrawlWebsite(req *CrawlWebsiteReq, res *CrawlWebsiteRes) error {
	log.Printf("Start crawling %s (depth= %d)\n", req.URL, req.Depth)
	domain := getDomainName(req.URL)
	// store domain name if hasn't already
	if _, in := savedDomains[domain]; !in {
		savedDomains[domain] = struct{}{}
	}
	if req.Depth == 0 {
		return nil
	}
	// get absolute path
	fmtRequestURL := filterAddress(req.URL, domain)
	links := getLinks(fmtRequestURL)
	log.Println("Done Crawling.")
	// could try to distribute this work
	for _, link := range links {
		if !contains(webGraph[fmtRequestURL], link) {
			webGraph[fmtRequestURL] = append(webGraph[fmtRequestURL], link)
		}
	}
	// doing this means each depth will return to server
	res.Depth--
	res.Links = links
	
	// for node, edges := range webGraph {
	// 	log.Printf("\nNode:\t%s\nEdges:\t%s\n", node, edges)
	// }


	return nil
}

// extract links from a web page
// returns a string of formatted link
// remove relative path of if a link doesn't start with http, add to it a domain name
func getLinks(uri string) (links []string) {
	resp, _ := http.Get(uri)
	defer resp.Body.Close()
	
	domain := getDomainName(uri)
	z := html.NewTokenizer(resp.Body)
	uniqueLinks := make(map[string]bool)

	for {
		tt := z.Next()
		switch tt {
			case html.ErrorToken:
				return
			case html.StartTagToken, html.EndTagToken:
				token := z.Token()
				if "a" == token.Data {
					// add only unique links
					for _, attr := range token.Attr {
						if attr.Key == "href" && ! uniqueLinks[attr.Val] {
							// check if the path is relative
							// if not append
							links = append(links, filterAddress(attr.Val, domain))
							uniqueLinks[attr.Val] = true
						}
					}
				}
		}
	}
	return
}

//*************************************************************
