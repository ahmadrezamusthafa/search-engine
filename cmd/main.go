package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/handler"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/router"
	"github.com/ahmadrezamusthafa/search-engine/pkg/badgerdb"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db := badgerdb.NewBadgerDB("./db")
	defer db.Close()

	searchEngine := engine.NewSearchEngine(cfg.BM25, db)

	h := handler.NewHandler(searchEngine)
	r := router.NewRouter(h)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s...", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
