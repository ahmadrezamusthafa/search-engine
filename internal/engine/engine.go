package engine

import (
	"encoding/json"
	"errors"
	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"github.com/dgraph-io/badger/v4"
	"log"
	"math"
	"sort"
	"sync"
)

type SearchEngine struct {
	mu           sync.RWMutex
	db           *badger.DB
	docTokensLen map[string]int
	termDocCount map[string]int
	tokenLen     int
	docCount     int
	k1           float64
	b            float64
}

func NewSearchEngine(config config.BM25Config, db *badger.DB) (*SearchEngine, error) {
	return &SearchEngine{
		docTokensLen: make(map[string]int),
		termDocCount: make(map[string]int),
		db:           db,
		k1:           config.K1,
		b:            config.B,
	}, nil
}

func (se *SearchEngine) StoreDocument(docID string, tokens []string, contents ...structs.Content) {
	se.mu.Lock()
	defer se.mu.Unlock()

	tokenFrequency := make(map[string]int)
	for _, token := range tokens {
		tokenFrequency[token]++
	}

	se.docTokensLen[docID] = len(tokens)
	se.tokenLen += len(tokens)
	se.docCount++

	for token, freq := range tokenFrequency {
		se.termDocCount[token]++

		err := se.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("index:" + token))
			if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
				return err
			}

			var currentData map[string]int
			if err == nil {
				err = item.Value(func(val []byte) error {
					return json.Unmarshal(val, &currentData)
				})
				if err != nil {
					return err
				}
			} else {
				currentData = make(map[string]int)
			}

			currentData[docID] = freq

			updatedData, err := json.Marshal(currentData)
			if err != nil {
				return err
			}

			err = txn.Set([]byte("index:"+token), updatedData)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}

	}

	if len(contents) > 0 {
		err := se.db.Update(func(txn *badger.Txn) error {
			val, err := json.Marshal(contents[0].Object)
			if err != nil {
				return err
			}
			err = txn.Set([]byte("data:"+docID), val)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (se *SearchEngine) Search(queries ...string) []structs.SearchResult {
	se.mu.RLock()
	defer se.mu.RUnlock()

	if len(queries) == 0 {
		return nil
	}

	avgDocLen := se.calculateAvgDocLength()
	docScores := make(map[string]float64)

	for _, query := range queries {
		err := se.db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("index:" + query))
			if err != nil {
				return err
			}

			var docFreqMap map[string]int
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &docFreqMap)
			})
			if err != nil {
				return err
			}

			for docID, tf := range docFreqMap {
				docLen := se.docTokensLen[docID]
				bm25Score := se.calculateBM25(tf, se.termDocCount[query], docLen, avgDocLen, se.k1, se.b)
				docScores[docID] += bm25Score
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}

	results := make([]structs.SearchResult, 0, len(docScores))
	for docID, score := range docScores {
		results = append(results, structs.SearchResult{ID: docID, Score: score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > 0 {
		results = util.GetTopItems(results, 3)
		err := se.db.View(func(txn *badger.Txn) error {
			for i, result := range results {
				item, err := txn.Get([]byte("data:" + result.ID))
				if err != nil {
					return err
				}
				var value map[string]interface{}
				err = item.Value(func(val []byte) error {
					return json.Unmarshal(val, &value)
				})
				if err != nil {
					return err
				}
				results[i].Data = value
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}

	return results
}

func (se *SearchEngine) calculateAvgDocLength() int {
	if se.docCount == 0 {
		return 0
	}
	return se.tokenLen / se.docCount
}

func (se *SearchEngine) calculateBM25(tf, df, docLen, avgDocLen int, k1, b float64) float64 {
	if df == 0 || avgDocLen == 0 {
		return 0
	}

	idf := math.Log((float64(se.docCount)-float64(df)+0.5)/(float64(df)+0.5) + 1)
	tfWeight := (float64(tf) * (k1 + 1)) / (float64(tf) + k1*(1-b+b*float64(docLen)/float64(avgDocLen)))
	return idf * tfWeight
}
