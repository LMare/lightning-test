package httpapi

import (
	"net/http"
)

func NewVersionHandler(version string) *VersionHandler {
	return &VersionHandler{version: version}
}

type VersionHandler struct {
	version string
}

// HandleVersion returns the current version of the application
func (h *VersionHandler) HandleVersion(w http.ResponseWriter, r *http.Request) {
	if IsHTMX(r) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<p>" + h.version + "</p>"))
	} else {
		JSONResponse(w, h.version)
	}
}
