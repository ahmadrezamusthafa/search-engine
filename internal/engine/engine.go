package engine

import (
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
}

func Search(queries ...string) []string {
	mu.RLock()
	defer mu.RUnlock()

	results := make([]string, 0)
	docMatches := make(map[string]int)

	for _, query := range queries {
		for docID, tokens := range index {
			if _, found := tokens[query]; found {
				if _, ok := docMatches[docID]; !ok {
					docMatches[docID] = 0
				}
				docMatches[docID]++
			}
		}
	}

	for docID, matchCount := range docMatches {
		if matchCount == len(queries) {
			results = append(results, docID)
		}
	}

	return results
}
