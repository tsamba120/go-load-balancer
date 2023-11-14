package main

import (
	"encoding/json"
	"fmt"
	"os"
	// "net/http"
)

type Server struct {
	Address           string
	ActiveConnections int
	Healthy           bool
}

type LoadBalancerConfig struct {
	Port                string   `json:"port"`
	HealthCheckInterval int      `json:"health_check_interval"`
	Servers             []Server `json:"servers"`
}

func main() {
	// read from config - either file i/o or user viper
	// instantiate into load balancer struct
	// lbConfig := getConfig()

	// start health check for each server in separate go routine

	// listen and serve
}

func performHealthCheck()

func getConfig() LoadBalancerConfig {
	data, err := os.ReadFile("./cmd/load_balancer/load_balancer_config.json")
	if err != nil {
		panic(err)
	}

	json_data := make(map[string]interface{})
	var lb LoadBalancerConfig

	json.Unmarshal([]byte(data), &json_data)

	lb.Port = json_data["port"].(string)
	lb.HealthCheckInterval = int(json_data["health_check_interval"].(float64))

	serverNameList := json_data["servers"].([]interface{})

	var serverList []Server

	for _, server := range serverNameList {
		s := Server{
			Address:           server.(string),
			ActiveConnections: 0,
			Healthy:           true,
		}
		serverList = append(serverList, s)
	}

	lb.Servers = serverList

	fmt.Printf("Load balancer port: %s\n", lb.Port)
	fmt.Printf("Health check interval (ms): %d\n", lb.HealthCheckInterval)
	fmt.Println(serverList)

	return lb

}
