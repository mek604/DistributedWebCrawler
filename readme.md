# Distributed Web Crawler
## CSC462, Summer 2020



### Usage:

#### Example (Run each command separately)
```
go run server.go :8000 127.0.0.1:4000 5
go run worker.go 127.0.0.1:8000 8001
go run worker.go 127.0.0.1:8000 80002
...
go run client.go -c 127.0.0.1:4000 http://www.facebook.com 3
```

#### For server.go
```
go run server.go [worker ip:port] [client ip:port] [# of samples to measure RTT]
go run server.go [accepting :port] [client ip:port] [# of samples to measure RTT]
```

#### For worker.go
```
go run worker.go [server ip:port] [worker port]
```

#### For client.go
```
GetWorkers:
go run client.go -g [server ip:port]

Crawl:
go run client.go -c [server ip:port] [url] [depth]

Domains:
go run client.go -d [server ip:port] [workerIP]

Overlap:
go run client.go -o [server ip:port] [url1] [url2]

EC1 (Page rank):
go run client.go -r [server ip:port] [url1]

EC2 (Search):
go run client.go -s [server ip:port] [string]

Example:
go run client.go -g 127.0.0.1:19001
```

## <u>Work Remains To Be Done</u>
+ Modify worker to continue crawling if depth has not reached
+ Add web page overlap computation

