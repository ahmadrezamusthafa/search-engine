package main

import (
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
	"github.com/ahmadrezamusthafa/search-engine/pkg/badgerdb"
	"log"
	"testing"
)

/*
PERFORMANCE
BenchmarkSearch-8   	   42193	     26462 ns/op
BenchmarkSearch-8   	 2919550	       403.5 ns/op
BenchmarkSearch-8   	 2946054	       390.2 ns/op - full using memory
BenchmarkSearch-8   	  523896	      3724 ns/op - optimize memory usage <CURRENT>
*/

func BenchmarkSearch(b *testing.B) {
	cfg, err := config.LoadConfig("../../config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db := badgerdb.NewBadgerDB("./db")
	defer db.Close()

	searchEngine, err := engine.NewSearchEngine(cfg.BM25, db)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}

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
