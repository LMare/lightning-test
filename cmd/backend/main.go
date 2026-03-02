package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Lmare/lightning-playground"
	exception "github.com/Lmare/lightning-playground/backend/exception"
)

func main() {
	cfg := config.Load()
	exception.ConfigureProjectBasePath(cfg.ProjectPath)
	db, err := initDB(cfg)
	if err != nil {
		//log.Fatal("Failed to initialize database:", err)
		fmt.Printf("Failed to initialize database: %v\n", err)
	}
	repos := initRepositories(db)
	factories := initFactories()
	services := initServices(repos, factories, cfg)
	handlers := initHandlers(services, cfg)
	router := initRouter(handlers)

	fmt.Printf("Server Backend started : %s:%s\n", cfg.BackendUrl, cfg.BackendPort)
	log.Fatal(http.ListenAndServe(":"+cfg.BackendPort, router))
}
