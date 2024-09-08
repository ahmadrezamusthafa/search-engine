package main

import (
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
	"github.com/ahmadrezamusthafa/search-engine/pkg/badgerdb"
	"github.com/ahmadrezamusthafa/search-engine/pkg/redisdb"
	"log"
	"testing"
)

/*
PERFORMANCE
BenchmarkBadgerSearchEngine-8   	   95324	     14439 ns/op <CURRENT>
*/

func BenchmarkBadgerSearchEngine(b *testing.B) {
	cfg, err := config.LoadConfig("../../config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db := badgerdb.NewBadgerDB(cfg.Badger)
	defer db.Close()

	searchEngine := engine.NewBadgerSearchEngine(cfg.BM25, db)

	searchEngine.StoreDocument("doc1", []string{"abc", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc2", []string{"bvbv", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3a", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3b", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3c", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc4", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri", "abc"})

	for n := 0; n < b.N; n++ {
		searchEngine.Search("abc")
	}
}

/*
PERFORMANCE
BenchmarkRedisSearchEngine-8    	    5712	    207930 ns/op <CURRENT>
*/
func BenchmarkRedisSearchEngine(b *testing.B) {
	cfg, err := config.LoadConfig("../../config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	redis := redisdb.NewRedis(cfg.Redis)
	defer redis.Close()

	searchEngine := engine.NewRedisSearchEngine(cfg.BM25, redis)

	searchEngine.StoreDocument("doc1", []string{"abc", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc2", []string{"bvbv", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3a", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3b", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc3c", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri"})
	searchEngine.StoreDocument("doc4", []string{"hgh", "nasbdm", "aksjdhaks", "iuyiuweyri", "abc"})

	for n := 0; n < b.N; n++ {
		searchEngine.Search("abc")
	}
}
