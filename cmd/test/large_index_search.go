package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
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

func generateLargeIndex(numDocs int) {
	for i := 0; i < numDocs; i++ {
		docID := fmt.Sprintf("doc-%d", i+1)
		tokens := generateDocument()
		engine.StoreDocument(docID, tokens)

		if (i+1)%1000 == 0 {
			fmt.Printf("Indexed %d documents\n", i+1)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numDocs := 100000
	start := time.Now()
	generateLargeIndex(numDocs)
	fmt.Printf("Indexed %d documents in %v\n", numDocs, time.Since(start))

	queries := []string{"abc"}
	start = time.Now()
	results := engine.Search(queries...)
	fmt.Printf("Search results for %v: %d documents found in %v\n", queries, len(results), time.Since(start))

	for _, result := range results {
		fmt.Println(result)
	}
}
