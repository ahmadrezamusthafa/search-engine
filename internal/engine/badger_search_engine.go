package engine

import (
	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"github.com/ahmadrezamusthafa/search-engine/pkg/badgerdb"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

type BadgerSearchEngine struct {
	mu       sync.RWMutex
	badgerDB *badgerdb.BadgerDB
	tokenLen int
	docCount int
	k1       float64
	b        float64
}

const BadgerTTL = 2 * time.Hour

func NewBadgerSearchEngine(config config.BM25Config, badgerDB *badgerdb.BadgerDB) ISearchEngine {
	tokenLen, docCount := repopulateDataFromBadger(badgerDB)
	return &BadgerSearchEngine{
		tokenLen: tokenLen,
		docCount: docCount,
		badgerDB: badgerDB,
		k1:       config.K1,
		b:        config.B,
	}
}

func repopulateDataFromBadger(badgerDB *badgerdb.BadgerDB) (int, int) {
	tokenLen, err := badgerDB.GetInt("tokenLen")
	if err != nil {
		return 0, 0
	}
	docCount, err := badgerDB.GetInt("docCount")
	if err != nil {
		return 0, 0
	}
	return tokenLen, docCount
}

func (se *BadgerSearchEngine) StoreDocument(docID string, tokens []string, contents ...structs.Content) {
	se.mu.Lock()
	defer se.mu.Unlock()

	tokenFrequency := make(map[string]int)
	for _, token := range tokens {
		tokenFrequency[token]++
	}

	se.tokenLen += len(tokens)
	se.docCount++

	for token, freq := range tokenFrequency {
		termDocCount, err := se.badgerDB.GetInt("termDocCount:" + token)
		if err != nil {
			log.Println(err)
		}
		termDocCount++
		err = se.badgerDB.SetInt("termDocCount:"+token, termDocCount, BadgerTTL)
		if err != nil {
			log.Println(err)
		}
		var currentIndexData map[string]int
		err = se.badgerDB.GetObject("index:"+token, &currentIndexData)
		if err != nil {
			log.Println(err)
		}
		if currentIndexData == nil {
			currentIndexData = make(map[string]int)
		}
		currentIndexData[docID] = freq
		err = se.badgerDB.SetObject("index:"+token, currentIndexData, BadgerTTL)
		if err != nil {
			log.Println(err)
		}
	}

	err := se.badgerDB.SetIntegers(BadgerTTL,
		badgerdb.KVInt{Key: "docTokensLen:" + docID, Value: len(tokens)},
		badgerdb.KVInt{Key: "tokenLen", Value: se.tokenLen},
		badgerdb.KVInt{Key: "docCount", Value: se.docCount},
	)
	if err != nil {
		log.Println(err)
	}

	if len(contents) > 0 {
		data := map[string]interface{}{
			"string": contents[0].String,
			"object": contents[0].Object,
		}
		err := se.badgerDB.SetObject("data:"+docID, data, BadgerTTL)
		if err != nil {
			log.Println(err)
		}
	}
}

func (se *BadgerSearchEngine) Search(queries ...string) []structs.SearchResult {
	se.mu.RLock()
	defer se.mu.RUnlock()

	if len(queries) == 0 {
		return nil
	}

	avgDocLen := se.calculateAvgDocLength()
	docScores := make(map[string]float64)

	for _, query := range queries {
		var docFreqMap map[string]int
		err := se.badgerDB.GetObject("index:"+query, &docFreqMap)
		if err != nil {
			return nil
		}

		if docFreqMap == nil {
			continue
		}

		for docID, tf := range docFreqMap {
			docLen, err := se.badgerDB.GetInt("docTokensLen:" + docID)
			if err != nil {
				log.Println(err)
			}

			termDocCount, err := se.badgerDB.GetInt("termDocCount:" + query)
			if err != nil {
				log.Println(err)
			}

			bm25Score := se.calculateBM25(tf, termDocCount, docLen, avgDocLen, se.k1, se.b)
			docScores[docID] += bm25Score
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
		for i, result := range results {
			var value map[string]interface{}
			err := se.badgerDB.GetObject("data:"+result.ID, &value)
			if err != nil {
				log.Println(err)
				continue
			}
			results[i].Data = value
		}
	}

	return results
}

func (se *BadgerSearchEngine) GetPersistenceType() string {
	return "BadgerDB"
}

func (se *BadgerSearchEngine) calculateAvgDocLength() int {
	if se.docCount == 0 {
		return 0
	}
	return se.tokenLen / se.docCount
}

func (se *BadgerSearchEngine) calculateBM25(tf, df, docLen, avgDocLen int, k1, b float64) float64 {
	if df == 0 || avgDocLen == 0 {
		return 0
	}

	idf := math.Log((float64(se.docCount)-float64(df)+0.5)/(float64(df)+0.5) + 1)
	tfWeight := (float64(tf) * (k1 + 1)) / (float64(tf) + k1*(1-b+b*float64(docLen)/float64(avgDocLen)))
	return idf * tfWeight
}
