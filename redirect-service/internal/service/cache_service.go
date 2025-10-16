package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	RefreshTTL(ctx context.Context, key string, ttl time.Duration) error
}

type cacheService struct {
	redis  *database.Redis
	config *config.Config
}

func NewCacheService(redis *database.Redis, cfg *config.Config) CacheService {
	return &cacheService{
		redis:  redis,
		config: cfg,
	}
}

func (s *cacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := s.redis.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss: key %s not found", key)
		}
		return fmt.Errorf("redis get error: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

func (s *cacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := s.redis.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	log.Printf("✅ Cached key: %s for %v", key, ttl)
	return nil
}

func (s *cacheService) Delete(ctx context.Context, key string) error {
	if err := s.redis.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}

	log.Printf("🗑️  Deleted cache key: %s", key)
	return nil
}

func (s *cacheService) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := s.redis.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}

	return exists > 0, nil
}

func (s *cacheService) RefreshTTL(ctx context.Context, key string, ttl time.Duration) error {
	if err := s.redis.Client.Expire(ctx, key, ttl).Err(); err != nil {
		log.Printf("⚠️  Failed to refresh cache TTL for key %s: %v", key, err)
		return fmt.Errorf("redis expire error: %w", err)
	}

	return nil
}
