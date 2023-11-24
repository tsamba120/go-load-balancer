package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

type Server struct {
	URL          	  *url.URL 
	ActiveConnections int
	HealthEndpoint    string
	Healthy           bool
	Mutex			  *sync.Mutex
}

type LoadBalancerConfig struct {
	Port                string   `json:"port"`
	HealthEndpoint      string   `json:"health_endpoint"`
	HealthCheckInterval string   `json:"health_check_interval"`
	Servers             []*Server `json:"servers"`
}

func main() {
	// read from config - either file i/o or user viper
	// instantiate into load balancer struct
	lbConfig := getConfig()

	healthCheckInterval, err := time.ParseDuration(lbConfig.HealthCheckInterval)
	if err != nil {
		panic(err)
	}

	// start health check for each server in separate go routine
	for _, server := range lbConfig.Servers {
		// go routine to run a health check for each server
		go func(s *Server) {
			for range time.Tick(healthCheckInterval) {
				performHealthCheck(s)
			}
		}(server)
	}


	// create mux and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/get", lbConfig.serveGetRequest)
	mux.HandleFunc("/put", lbConfig.servePutRequest)

	// instantiate server and listen
	srv := &http.Server{
		Addr: 	":4040",
		Handler: mux,
	}

	fmt.Println("Starting load balancer on :4040")
	err = srv.ListenAndServe()
	panic(err)
}

func (server Server) Proxy(requestPath string) *httputil.ReverseProxy {
	if requestPath != "put" && requestPath != "get" {
		// TODO FAIL GRACEFULLY - don't kill the webserver
		panic("request path needs to be either 'put' or 'get")
	}
	baseUrl := *server.URL
	baseUrl.Path = path.Join(baseUrl.Path, "data", requestPath)
	return httputil.NewSingleHostReverseProxy(&baseUrl)
}

func performHealthCheck(server *Server) {
	resp, err := http.Get(server.URL.String() + server.HealthEndpoint)
	if err != nil || resp.StatusCode != 200 {
		server.Healthy = false
	} else {
		server.Healthy = true
	}
}


func getConfig() LoadBalancerConfig {
	// TODO use cobra
	data, err := os.ReadFile("./cmd/load_balancer/load_balancer_config.json")
	if err != nil {
		panic(err)
	}

	json_data := make(map[string]interface{})
	var lb LoadBalancerConfig

	json.Unmarshal([]byte(data), &json_data)

	lb.Port = json_data["port"].(string)
	lb.HealthEndpoint = json_data["health_endpoint"].(string)
	lb.HealthCheckInterval = json_data["health_check_interval"].(string)

	serverNameList := json_data["servers"].([]interface{})

	var serverList []*Server

	for _, server := range serverNameList {
		serverURL, _ := url.Parse(server.(string))
		s := &Server{
			URL:           	   serverURL,
			Healthy:           true,
			Mutex:			   &sync.Mutex{},
		}
		serverList = append(serverList, s)
	}

	lb.Servers = serverList

	// fmt.Printf("Load balancer port: %s\n", lb.Port)
	// fmt.Printf("Health check endpoint: %s\n", lb.HealthEndpoint)
	// fmt.Printf("Health check interval (ms): %s\n", lb.HealthCheckInterval)
	// for _, server := range serverList {
	// 	fmt.Println(server)
	// }

	return lb

}
