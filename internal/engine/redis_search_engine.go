package engine

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
)

type RedisSearchEngine struct {
	mu       sync.RWMutex
	redisDB  *redis.Client
	ctx      context.Context
	tokenLen int
	docCount int
	k1       float64
	b        float64
}

const RedisTTL = 2 * time.Hour

func NewRedisSearchEngine(config config.BM25Config, redisDB *redis.Client) ISearchEngine {
	ctx := context.Background()
	tokenLen, docCount := repopulateDataFromRedis(ctx, redisDB)
	return &RedisSearchEngine{
		tokenLen: tokenLen,
		docCount: docCount,
		redisDB:  redisDB,
		ctx:      ctx,
		k1:       config.K1,
		b:        config.B,
	}
}

func repopulateDataFromRedis(ctx context.Context, redisDB *redis.Client) (int, int) {
	tokenLen, err := redisDB.Get(ctx, "tokenLen").Int()
	if err != nil {
		return 0, 0
	}
	docCount, err := redisDB.Get(ctx, "docCount").Int()
	if err != nil {
		return 0, 0
	}
	return tokenLen, docCount
}

func (se *RedisSearchEngine) StoreDocument(docID string, tokens []string, contents ...structs.Content) {
	se.mu.Lock()
	defer se.mu.Unlock()

	tokenFrequency := make(map[string]int)
	for _, token := range tokens {
		tokenFrequency[token]++
	}

	se.tokenLen += len(tokens)
	se.docCount++

	for token, freq := range tokenFrequency {
		termDocCount, err := se.redisDB.Get(se.ctx, "termDocCount:"+token).Int()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Println(err)
		}
		termDocCount++
		err = se.redisDB.Set(se.ctx, "termDocCount:"+token, termDocCount, RedisTTL).Err()
		if err != nil {
			log.Println(err)
		}

		docFreqMap := make(map[string]int)
		res, err := se.redisDB.Get(se.ctx, "index:"+token).Bytes()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Println(err)
		}

		if res != nil {
			err = json.Unmarshal(res, &docFreqMap)
			if err != nil {
				log.Println(err)
			}
		}

		docFreqMap[docID] = freq
		docFreqMapBytes, err := json.Marshal(docFreqMap)
		if err != nil {
			log.Println(err)
		}

		err = se.redisDB.Set(se.ctx, "index:"+token, docFreqMapBytes, RedisTTL).Err()
		if err != nil {
			log.Println(err)
		}
	}

	err := se.redisDB.Set(se.ctx, "docTokensLen:"+docID, len(tokens), RedisTTL).Err()
	if err != nil {
		log.Println(err)
	}
	err = se.redisDB.Set(se.ctx, "tokenLen", se.tokenLen, RedisTTL).Err()
	if err != nil {
		log.Println(err)
	}
	err = se.redisDB.Set(se.ctx, "docCount", se.docCount, RedisTTL).Err()
	if err != nil {
		log.Println(err)
	}

	if len(contents) > 0 {
		data := map[string]interface{}{
			"string": contents[0].String,
			"object": contents[0].Object,
		}
		contentBytes, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		}

		err = se.redisDB.Set(se.ctx, "data:"+docID, contentBytes, RedisTTL).Err()
		if err != nil {
			log.Println(err)
		}
	}
}

func (se *RedisSearchEngine) Search(queries ...string) []structs.SearchResult {
	se.mu.RLock()
	defer se.mu.RUnlock()

	if len(queries) == 0 {
		return nil
	}

	avgDocLen := se.calculateAvgDocLength()
	docScores := make(map[string]float64)

	for _, query := range queries {
		docFreqMap := make(map[string]int)
		res, err := se.redisDB.Get(se.ctx, "index:"+query).Bytes()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Println(err)
			return nil
		}

		if res != nil {
			err = json.Unmarshal(res, &docFreqMap)
			if err != nil {
				log.Println(err)
			}
		}

		if len(docFreqMap) == 0 {
			continue
		}

		for docID, tf := range docFreqMap {
			docLen, err := se.redisDB.Get(se.ctx, "docTokensLen:"+docID).Int()
			if err != nil && !errors.Is(err, redis.Nil) {
				log.Println(err)
			}

			termDocCount, err := se.redisDB.Get(se.ctx, "termDocCount:"+query).Int()
			if err != nil && !errors.Is(err, redis.Nil) {
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
			res, err := se.redisDB.Get(se.ctx, "data:"+result.ID).Bytes()
			if err != nil && !errors.Is(err, redis.Nil) {
				log.Println(err)
				continue
			}

			if res != nil {
				err = json.Unmarshal(res, &value)
				if err != nil {
					log.Println(err)
				}
			}
			results[i].Data = value
		}
	}

	return results
}

func (se *RedisSearchEngine) GetPersistenceType() string {
	return "Redis"
}

func (se *RedisSearchEngine) calculateAvgDocLength() int {
	if se.docCount == 0 {
		return 0
	}
	return se.tokenLen / se.docCount
}

func (se *RedisSearchEngine) calculateBM25(tf, df, docLen, avgDocLen int, k1, b float64) float64 {
	if df == 0 || avgDocLen == 0 {
		return 0
	}

	idf := math.Log((float64(se.docCount)-float64(df)+0.5)/(float64(df)+0.5) + 1)
	tfWeight := (float64(tf) * (k1 + 1)) / (float64(tf) + k1*(1-b+b*float64(docLen)/float64(avgDocLen)))
	return idf * tfWeight
}
