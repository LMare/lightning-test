package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"

	config "github.com/Lmare/lightning-test"
)

func main() {
	startProxy()
	startServer()
}

// Start the static File server and template server
func startServer() {
	cfg := config.Load()

	// Serveur de fichiers statiques
	staticPath := filepath.Join(".", "frontend", "asset")
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/asset/", http.StripPrefix("/asset/", fs))

	// Routes frontend
 	http.HandleFunc("/", pageHandler)


 	log.Printf("Frontend sur %s:%s", cfg.FrontendUrl, cfg.FrontendPort)
	log.Fatal(http.ListenAndServe(":" + cfg.FrontendPort, nil))

}

// Reverse proxy pour /api/*
// Remplace certaines variables du backend dans la réponse
func startProxy() {
	cfg := config.Load()

	// Cible du backend
	target, err := url.Parse(cfg.BackendUrl + ":" + cfg.BackendPort)
	if err != nil {
		log.Fatal(err)
	}

	// Reverse proxy pour /api/*
	proxy := httputil.NewSingleHostReverseProxy(target)
	// Interception et modification de la réponse
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Ne pas toucher aux flux SSE
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/event-stream") {
			return nil
		}

		// Lire et modifier le corps pour les autres types
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		// Modifier le corps (ex: remplacer un placeholder)
		modifiedBody := bytes.ReplaceAll(bodyBytes, []byte("__API_PREFIX__"), []byte("/api"))

		// Réinjecter le corps modifié
		resp.Body = io.NopCloser(bytes.NewReader(modifiedBody))
		resp.ContentLength = int64(len(modifiedBody))
		resp.Header.Set("Content-Length", fmt.Sprint(len(modifiedBody)))

		return nil
	}
	// interception de la requête
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// Retire le préfixe "/api" du chemin
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Serveur proxy lancé sur %s:%s\n", cfg.FrontendUrl, cfg.FrontendPort)
}



// render fragment html ou page avec le layout
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
			log.Printf("error : %v", err)
			http.NotFound(w, r)
		}
		return
	}

	// Sinon → layout complet avec fragment inclus
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, page+".html", nil)
	if err != nil {
		log.Printf("error : %v", err)
		http.NotFound(w, r)
		return
	}

	err = templates.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title":   strings.Title(page),
		"ActivePage": page,
		"Content": template.HTML(buf.String()), // contenu déjà rendu,
	})
	if err != nil {
		log.Printf("error : %v", err)
		http.NotFound(w, r)
	}


}
