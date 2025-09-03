package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Lmare/lightning-test/internal/handler"
)

const server = "http://localhost"
const port = ":8080"

func main() {
	startServer()
}

func startServer() {
	handler.Init()
	for _, route := range handler.Routes {
		http.HandleFunc(route.Path, route.Callback)
	}

	fmt.Printf("Server Backend started : %s%s\n", server, port)
	log.Fatal(http.ListenAndServe(port, nil))
}
