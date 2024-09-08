package engine

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"github.com/ahmadrezamusthafa/search-engine/pkg/badgerdb"
	"github.com/go-redis/redis/v8"
)

type ISearchEngine interface {
	StoreDocument(docID string, tokens []string, contents ...structs.Content)
	Search(queries ...string) []structs.SearchResult
	GetPersistenceType() string
}

const (
	PersistenceRedis  string = "redis"
	PersistenceBadger string = "badger"
)

func NewSearchEngine[T any](
	persistenceType string,
	cfg *config.Config,
	db T) (ISearchEngine, error) {

	switch persistenceType {
	case "redis":
		if redisClient, ok := any(db).(*redis.Client); ok {
			return NewRedisSearchEngine(cfg.BM25, redisClient), nil
		}
		return nil, fmt.Errorf("invalid type for Redis persistence")
	case "badger":
		if badgerDB, ok := any(db).(*badgerdb.BadgerDB); ok {
			return NewBadgerSearchEngine(cfg.BM25, badgerDB), nil
		}
		return nil, fmt.Errorf("invalid type for BadgerDB persistence")
	default:
		return nil, fmt.Errorf("unsupported persistence type: use redis or badger as search engine persistence")
	}
}
