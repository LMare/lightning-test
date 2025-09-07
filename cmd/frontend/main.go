package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"

	config "github.com/Lmare/lightning-test"
)

func main() {
	startServer()
}

// Start the static File server and the proxy
func startServer() {
	cfg := config.Load()
	// Cible du backend
	target, err := url.Parse(cfg.BackendUrl + ":" + cfg.BackendPort)
	if err != nil {
		log.Fatal(err)
	}

	// Reverse proxy pour /api/*
	proxy := httputil.NewSingleHostReverseProxy(target)
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// Retire le préfixe "/api" du chemin
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		proxy.ServeHTTP(w, r)
	})

	// Serveur de fichiers statiques
	staticDir := filepath.Join(".", "frontend/static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)

	log.Printf("Serveur proxy lancé sur %s:%s\n", cfg.FrontendUrl, cfg.FrontendPort)
	log.Printf("→ /api/* redirigé vers  %s:%s\n", cfg.BackendUrl, cfg.BackendPort)
	log.Println("→ / sert les fichiers depuis .frontend/static")
	http.ListenAndServe(":"+cfg.FrontendPort, nil)

}
