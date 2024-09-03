package storage

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	mu        sync.RWMutex
	index     = make(map[string]map[string]int) // docID -> token -> frequency
	docTokens = make(map[string][]string)       // docID -> tokens
)

func StoreDocument(docID string, tokens []string) {
	mu.Lock()
	defer mu.Unlock()

	tokenFrequency := make(map[string]int)
	for _, token := range tokens {
		tokenFrequency[token]++
	}
	index[docID] = tokenFrequency
	docTokens[docID] = tokens

	b, _ := json.MarshalIndent(index, "", "  ")
	fmt.Println("Index:", string(b))
	b, _ = json.MarshalIndent(docTokens, "", "  ")
	fmt.Println("Tokens:", string(b))
}

func Search(query string) []string {
	mu.RLock()
	defer mu.RUnlock()

	results := make([]string, 0)
	for docID, tokens := range index {
		if _, found := tokens[query]; found {
			results = append(results, docID)
		}
	}
	return results
}
