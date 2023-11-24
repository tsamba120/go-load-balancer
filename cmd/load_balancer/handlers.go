package main

import (
	"fmt"
	"net/http"
)

// Get the server from (slice of servers) with the smallest number of 
// active connections
func (lb LoadBalancerConfig) nextLeastActiveServer() *Server {
	leastActiveConnection := -1
	leastActiveServer := lb.Servers[0]
	for _, server := range lb.Servers {
		server.Mutex.Lock()
		if (server.ActiveConnections < leastActiveConnection || leastActiveConnection == -1) {
			leastActiveConnection = server.ActiveConnections
			leastActiveServer = server	
		}
		server.Mutex.Unlock()
	}
	return leastActiveServer
}


func (lb LoadBalancerConfig) handleRequest(w http.ResponseWriter, r *http.Request, requestType string) {
	server := lb.nextLeastActiveServer()

	server.Mutex.Lock()
	server.ActiveConnections++
	server.Mutex.Unlock()

	fmt.Println("Forwarding request now")

	// forward request
	server.Proxy(requestType).ServeHTTP(w, r)

	server.Mutex.Lock()
	server.ActiveConnections--
	server.Mutex.Unlock()

}


func (lb LoadBalancerConfig) servePutRequest(w http.ResponseWriter, r *http.Request) {
	lb.handleRequest(w, r, "put")
}


func (lb LoadBalancerConfig) serveGetRequest(w http.ResponseWriter, r *http.Request) {
	lb.handleRequest(w, r, "get")
}