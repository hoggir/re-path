package database

import (
	"context"
	"fmt"
	"log"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(cfg *config.Config) (*Redis, error) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: 10,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("✅ Redis connected successfully at %s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return &Redis{
		Client: client,
	}, nil
}

func (r *Redis) Close() error {
	log.Println("🔌 Closing Redis connection...")
	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("failed to close redis connection: %w", err)
	}

	log.Println("✅ Redis connection closed successfully")
	return nil
}
