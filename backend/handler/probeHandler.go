package handler

import (
	"net/http"

	Config "github.com/Lmare/lightning-playground"
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

// handleVersion returns the current version of the application
func handleVersion(w http.ResponseWriter, r *http.Request) {
	if IsHTMX(r) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<p>" + Config.Load().Version + "</p>"))
	} else {
		jsonResponse(w, Config.Load().Version)
	}
}

// handleMetrics returns the metrics of the application
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	// here you can add your metrics collection and return them in the response
	// for example, you can use Prometheus client library to collect and expose metrics
	jsonResponse(w, "metrics not implemented yet")
}
