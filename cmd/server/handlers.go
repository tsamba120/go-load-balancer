package main

import (
	// "fmt"
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Health string `json:"health"`
}

type putRequestBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type getResponseBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("/health")
	response := healthResponse{Health: "Healthy"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (app *application) put(w http.ResponseWriter, r *http.Request) {

	var b putRequestBody

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Could not unmarshal request body", http.StatusBadRequest)
	}
	app.logger.Printf("/data/put with body: '%v'\n", b)

	app.mu.Lock()
	(*app.keyValueStore)[b.Key] = b.Value
	app.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)

}

func (app *application) get(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("/data/get")
	key := r.URL.Query().Get("key")

	app.mu.Lock()
	value := (*app.keyValueStore)[key]
	app.mu.Unlock()

	response := getResponseBody{
		Key:   key,
		Value: value,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (app *application) getRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", app.health)
	mux.HandleFunc("/data/put", app.put)
	mux.HandleFunc("/data/get", app.get)

	return mux
}
