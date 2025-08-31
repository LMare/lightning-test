package main

import (
	"fmt"
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

	fmt.Printf("Server started : %s%s\n", server, port)
	http.ListenAndServe(port, nil)
}
