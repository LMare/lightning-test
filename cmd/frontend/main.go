package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

const server = "http://localhost"
const port = ":3000"

func main() {
	startServer()
}

func startServer() {
	// Cible du backend
	target, err := url.Parse("http://localhost:8080")
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
	staticDir := filepath.Join(".", "static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)

	log.Println("Serveur proxy lancé sur http://localhost:3000")
	log.Println("→ /api/* redirigé vers http://localhost:8080")
	log.Println("→ / sert les fichiers depuis ./static")
	http.ListenAndServe(":3000", nil)

}
