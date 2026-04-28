package database

import (
	"context"
	"log"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) *redis.Client {
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	ctx := context.Background()

	redisClient := redis.NewClient(opts)

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return redisClient
}
