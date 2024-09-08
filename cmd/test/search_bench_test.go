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
BenchmarkSearch-8   	  523896	      3724 ns/op - optimize memory usage
BenchmarkSearch-8   	   53751	     18972 ns/op - full using badger <CURRENT>
*/

func BenchmarkSearch(b *testing.B) {
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
