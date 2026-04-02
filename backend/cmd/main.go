package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Lmare/lightning-playground"
	"github.com/Lmare/lightning-playground/backend/internal/bootstrap"
	exception "github.com/Lmare/lightning-playground/backend/internal/shared/exception"
)

func main() {
	cfg := config.Load()
	exception.ConfigureProjectBasePath(cfg.ProjectPath)
	db, router, err := bootstrap.InitApp(cfg)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	executeMigrations(db)

	fmt.Printf("Server Backend started : %s:%s\n", cfg.BackendUrl, cfg.BackendPort)
	log.Fatal(http.ListenAndServe(":"+cfg.BackendPort, router))
}
