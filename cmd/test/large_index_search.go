package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
	"github.com/dgraph-io/badger/v4"
	"log"
	"math/rand"
	"time"
)

func generateRandomWord() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	word := make([]rune, rand.Intn(8)+3)
	for i := range word {
		word[i] = letters[rand.Intn(len(letters))]
	}
	return string(word)
}

func generateDocument() []string {
	numWords := rand.Intn(100) + 50
	words := make([]string, numWords)
	for i := range words {
		words[i] = generateRandomWord()
	}
	return words
}

func generateLargeIndex(searchEngine *engine.SearchEngine, numDocs int) {
	for i := 0; i < numDocs; i++ {
		docID := fmt.Sprintf("doc-%d", i+1)
		tokens := generateDocument()
		searchEngine.StoreDocument(docID, tokens)

		if (i+1)%1000 == 0 {
			fmt.Printf("Indexed %d documents\n", i+1)
		}
	}
}

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

	rand.Seed(time.Now().UnixNano())

	numDocs := 100000
	start := time.Now()
	generateLargeIndex(searchEngine, numDocs)
	fmt.Printf("Indexed %d documents in %v\n", numDocs, time.Since(start))

	queries := []string{"abc"}
	start = time.Now()
	results := searchEngine.Search(queries...)
	fmt.Printf("Search results for %v: %d documents found in %v\n", queries, len(results), time.Since(start))

	for _, result := range results {
		fmt.Println(result)
	}
}