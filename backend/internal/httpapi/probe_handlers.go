package httpapi

import (
	"net/http"
)

type ProbeHandler struct {
}

func NewProbeHandler() *ProbeHandler {
	return &ProbeHandler{}
}

// HandleHealth checks if the application is healthy and can respond to requests
func (h *ProbeHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, "OK")
}

// HandleReady checks if the application is ready to serve requests
func (h *ProbeHandler) HandleReady(w http.ResponseWriter, r *http.Request) {
	// here you can add checks to ensure that the application is ready to serve requests,
	// such as checking database connections, external services, etc.
	JSONResponse(w, "OK")
}
