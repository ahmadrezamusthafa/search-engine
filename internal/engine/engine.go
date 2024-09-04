package engine

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"log"
	"sync"
)

type SearchEngine struct {
	mu    sync.RWMutex
	index map[string]map[string]int
	db    *badger.DB
}

func NewSearchEngine(storageDir string) (*SearchEngine, error) {
	opts := badger.DefaultOptions(storageDir).WithLoggingLevel(badger.INFO)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &SearchEngine{
		index: make(map[string]map[string]int),
		db:    db,
	}, nil
}

func (se *SearchEngine) Close() error {
	return se.db.Close()
}

func (se *SearchEngine) StoreDocument(docID string, tokens []string) {
	se.mu.Lock()
	defer se.mu.Unlock()

	tokenFrequency := make(map[string]int)
	for _, token := range tokens {
		tokenFrequency[token]++
	}
	se.index[docID] = tokenFrequency

	data, err := json.Marshal(tokenFrequency)
	if err != nil {
		log.Printf("Failed to serialize document %s: %v", docID, err)
		return
	}

	err = se.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(docID), data)
	})
	if err != nil {
		log.Printf("Failed to store document %s in DB: %v", docID, err)
	}
}

func (se *SearchEngine) Search(queries ...string) []string {
	se.mu.RLock()
	defer se.mu.RUnlock()

	results := make([]string, 0)

	for docID, tokenMap := range se.index {
		tokens, err := se.getDocumentTokens(docID, tokenMap)
		if err != nil {
			log.Printf("Failed to retrieve document tokens for %s: %v", docID, err)
			continue
		}

		matches := true
		for _, query := range queries {
			if _, found := tokens[query]; !found {
				matches = false
				break
			}
		}
		if matches {
			results = append(results, docID)
		}
	}

	return results
}

func (se *SearchEngine) getDocumentTokens(docID string, tokenMap map[string]int) (map[string]int, error) {
	if tokenMap != nil {
		return tokenMap, nil
	}

	var tokens map[string]int
	err := se.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(docID))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &tokens)
		})
	})

	if err != nil {
		return nil, err
	}

	se.index[docID] = tokens
	return tokens, nil
}
