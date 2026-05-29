package database

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"cmdb-api/config"
)

var (
	Redis    *redis.Client
	Ctx      = context.Background()
	redisOnce sync.Once
)

func InitRedis(cfg *config.Config) {
	redisOnce.Do(func() {
		Redis = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
	})
}
