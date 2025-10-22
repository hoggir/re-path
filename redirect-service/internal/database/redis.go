package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
	logger logger.Logger
}

func NewRedis(cfg *config.Config, log logger.Logger) (*Redis, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Redis.ConnTimeout)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.ConnTimeout,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info("Redis connected successfully",
		"host", cfg.Redis.Host,
		"port", cfg.Redis.Port,
		"poolSize", cfg.Redis.PoolSize,
		"minIdleConns", cfg.Redis.MinIdleConns)

	return &Redis{
		Client: client,
		logger: log,
	}, nil
}

func (r *Redis) Close() error {
	r.logger.Info("closing Redis connection")
	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("failed to close redis connection: %w", err)
	}

	r.logger.Info("Redis connection closed successfully")
	return nil
}
