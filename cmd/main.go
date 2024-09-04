package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/handler"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/router"
	"github.com/dgraph-io/badger/v4"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	opts := badger.DefaultOptions("./db").WithLoggingLevel(badger.INFO)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.Close()

	searchEngine, err := engine.NewSearchEngine(cfg.BM25, db)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}

	h := handler.NewHandler(searchEngine)
	r := router.NewRouter(h)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s...", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
