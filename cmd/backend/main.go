package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Lmare/lightning-playground"
	app "github.com/Lmare/lightning-playground/backend/app"
	exception "github.com/Lmare/lightning-playground/backend/exception"
)

func main() {
	cfg := config.Load()
	exception.ConfigureProjectBasePath(cfg.ProjectPath)
	db, router, err := app.InitApp(cfg)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	executeMigrations(db)

	fmt.Printf("Server Backend started : %s:%s\n", cfg.BackendUrl, cfg.BackendPort)
	log.Fatal(http.ListenAndServe(":"+cfg.BackendPort, router))
}
