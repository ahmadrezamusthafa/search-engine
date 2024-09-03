package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"log"
	"net/http"

	"github.com/ahmadrezamusthafa/search-engine/internal/router"
)

func main() {
	r := router.NewRouter()

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
