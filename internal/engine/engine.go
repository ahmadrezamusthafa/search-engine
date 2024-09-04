package engine

import (
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"math"
	"sort"
	"sync"
)

type SearchEngine struct {
	mu           sync.RWMutex
	index        map[string]map[string]int
	docTokens    map[string][]string
	termDocCount map[string]int
	docCount     int

	// constant for BM25
	k1 float64
	b  float64
}

func NewSearchEngine(config config.BM25Config) (*SearchEngine, error) {
	return &SearchEngine{
		index:        make(map[string]map[string]int),
		docTokens:    make(map[string][]string),
		termDocCount: make(map[string]int),
		k1:           config.K1,
		b:            config.B,
	}, nil
}

func (se *SearchEngine) StoreDocument(docID string, tokens []string) {
	se.mu.Lock()
	defer se.mu.Unlock()

	tokenFrequency := make(map[string]int)
	for _, token := range tokens {
		tokenFrequency[token]++
	}

	se.docTokens[docID] = tokens
	se.docCount++

	for token, freq := range tokenFrequency {
		se.termDocCount[token]++
		if se.index[token] == nil {
			se.index[token] = make(map[string]int)
		}
		se.index[token][docID] = freq
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
		if docFreqMap, found := se.index[query]; found {
			for docID, tf := range docFreqMap {
				docLen := len(se.docTokens[docID])
				bm25Score := se.calculateBM25(tf, se.termDocCount[query], docLen, avgDocLen, se.k1, se.b)
				docScores[docID] += bm25Score
			}
		}
	}

	results := make([]structs.SearchResult, 0, len(docScores))
	for docID, score := range docScores {
		results = append(results, structs.SearchResult{ID: docID, Score: score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func (se *SearchEngine) calculateAvgDocLength() int {
	totalLength := 0
	for _, tokens := range se.docTokens {
		totalLength += len(tokens)
	}
	if se.docCount == 0 {
		return 0
	}
	return totalLength / se.docCount
}

func (se *SearchEngine) calculateBM25(tf, df, docLen, avgDocLen int, k1, b float64) float64 {
	idf := math.Log((float64(se.docCount)-float64(df)+0.5)/(float64(df)+0.5) + 1)
	tfWeight := (float64(tf) * (k1 + 1)) / (float64(tf) + k1*(1-b+b*float64(docLen)/float64(avgDocLen)))
	return idf * tfWeight
}
