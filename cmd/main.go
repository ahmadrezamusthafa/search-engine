package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/handler"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/router"
	"log"
	"net/http"
)

func main() {
	searchEngine, err := engine.NewSearchEngine("./badgerdb")
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}
	defer searchEngine.Close()

	h := handler.NewHandler(searchEngine)
	r := router.NewRouter(h)

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s...", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
