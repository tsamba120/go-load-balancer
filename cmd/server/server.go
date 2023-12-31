package main

import (
	"log"
	"net/http"
	"os"
	"sync"
)

type application struct {
	logger        *log.Logger
	keyValueStore *map[string]any
	mu            *sync.Mutex
}

func main() {
	addr := ":8080"
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	kvStore := make(map[string]interface{})

	app := application{
		logger:        infoLog,
		keyValueStore: &kvStore,
		mu:            &sync.Mutex{},
	}

	server := &http.Server{
		Addr:    addr,
		Handler: app.getRoutes(),
	}

	// TODO: log when running
	infoLog.Printf("Running server on %s\n", addr)
	infoLog.Fatal(server.ListenAndServe())

}
