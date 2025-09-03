package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Lmare/lightning-test"
	"github.com/Lmare/lightning-test/internal/handler"
)

func main() {
	startServer()
}

func startServer() {
	cfg := config.Load()
	handler.Init()
	for _, route := range handler.Routes {
		http.HandleFunc(route.Path, route.Callback)
	}

	fmt.Printf("Server Backend started : %s:%s\n", cfg.BackendUrl, cfg.BackendPort)
	log.Fatal(http.ListenAndServe(":"+cfg.BackendPort, nil))
}
