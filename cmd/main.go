package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/handler"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/router"
	"github.com/ahmadrezamusthafa/search-engine/pkg/badgerdb"
	"github.com/ahmadrezamusthafa/search-engine/pkg/redisdb"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	badgerDB := badgerdb.NewBadgerDB(cfg.Badger)
	defer badgerDB.Close()
	redis := redisdb.NewRedis(cfg.Redis)
	defer redis.Close()

	//Sample: Redis
	//searchEngine, err := engine.NewSearchEngine(engine.PersistenceRedis, cfg, redis)

	searchEngine, err := engine.NewSearchEngine(engine.PersistenceBadger, cfg, badgerDB)
	if err != nil {
		log.Fatalf("Error initiate search engine: %v", err)
	}

	h := handler.NewHandler(searchEngine)
	r := router.NewRouter(h)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s...", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
