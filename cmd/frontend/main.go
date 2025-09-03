package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

const server = "http://localhost"
const port = ":3000"

func main() {
	startServer()
}

func startServer() {
	staticDir, err := filepath.Abs(filepath.Join(".", "static"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server Frontend started : %s%s\n", server, port)
	log.Fatal(http.ListenAndServe(port, http.FileServer(http.Dir(staticDir))))
}
