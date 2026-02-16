package handler

import (
	"net/http"
)

// handleHealth checks if the application is healthy and can respond to requests
func handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, "OK")
}

// handleReady checks if the application is ready to serve requests
func handleReady(w http.ResponseWriter, r *http.Request) {
	// here you can add checks to ensure that the application is ready to serve requests,
	// such as checking database connections, external services, etc.
	jsonResponse(w, "OK")
}
