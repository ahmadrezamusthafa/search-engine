package redisdb

import (
	"fmt"
	"github.com/ahmadrezamusthafa/search-engine/config"
	"github.com/go-redis/redis/v8"
)

func NewRedis(cfg config.RedisConfig) *redis.Client {
	redisOptions := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:           cfg.DB,
		DialTimeout:  cfg.DialConnectTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.MaxActive,
		MinIdleConns: cfg.MaxIdle,
		MaxConnAge:   cfg.MaxConnLifetime,
	}

	if cfg.Password != "" {
		redisOptions.Password = cfg.Password
	}

	return redis.NewClient(redisOptions)
}
