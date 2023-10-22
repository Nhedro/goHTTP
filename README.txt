This repository builds a server that lsitens to a port and returns the number of requests in the last 60 seconds
In this repository there are 2 files that can be run through command line:

Server 
	It expects a port to listen and a name of file where to fetch and keep the timestamps of requests received
	"go run goHTTPServer/main.go <port> <filename>"
 
Client
	It is expected to receive a url(http://127.0.0.1:<port> if local),. number of requests to perform and delay between them
	"go run goHTTPClient/main.go <url> <numRequests> <delay>"