package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"cmdb-api/config"
)

var Redis *redis.Client
var Ctx = context.Background()

func InitRedis(cfg *config.Config) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}
