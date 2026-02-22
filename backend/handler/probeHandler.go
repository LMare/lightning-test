package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

// handleMetrics exposes Prometheus metrics
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
