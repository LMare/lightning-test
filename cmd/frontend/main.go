package main

import (
	"bytes"
	"html/template"
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

	staticPath := filepath.Join(".", "frontend", "asset")
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/asset/", http.StripPrefix("/asset/", fs))


	// Routes frontend
 	http.HandleFunc("/", pageHandler)

	log.Printf("Serveur proxy lancé sur %s:%s\n", cfg.FrontendUrl, cfg.FrontendPort)
 	log.Printf("Frontend sur %s:%s", cfg.FrontendUrl, cfg.FrontendPort)
	log.Fatal(http.ListenAndServe(":"+cfg.FrontendPort, nil))

}

func pageHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	// Charge tous les templates
	templates := template.Must(template.ParseGlob("frontend/templates/*.html"))


    page := strings.Trim(r.URL.Path, "/")
    if page == "" {
        page = "hello" // page par défaut
    }

    // Si requête HTMX → fragment seul
    if r.Header.Get("HX-Request") == "true" {
        err := templates.ExecuteTemplate(w, page+".html", nil)
        if err != nil {
            http.NotFound(w, r)
        }
        return
    }

    // Sinon → layout complet avec fragment inclus
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, page+".html", nil)
	if err != nil {
	    http.NotFound(w, r)
	    return
	}

    err = templates.ExecuteTemplate(w, "layout.html", map[string]any{
        "Title":   strings.Title(page),
		"ActivePage": page,
        "Content": template.HTML(buf.String()), // contenu déjà rendu,
    })
    if err != nil {
        http.NotFound(w, r)
    }


}
