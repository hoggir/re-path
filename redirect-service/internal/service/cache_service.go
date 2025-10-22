package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	RefreshTTL(ctx context.Context, key string, ttl time.Duration) error
	SetInvalidationFlag(ctx context.Context, key string, ttl time.Duration) error
}

type cacheService struct {
	redis  *database.Redis
	config *config.Config
	logger logger.Logger
}

func NewCacheService(redis *database.Redis, cfg *config.Config, log logger.Logger) CacheService {
	return &cacheService{
		redis:  redis,
		config: cfg,
		logger: log,
	}
}

func (s *cacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := s.redis.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss is not really an error, just return a specific error
			return fmt.Errorf("cache miss: key %s not found", key)
		}
		return domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "Get").
			Wrap(err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "Unmarshal").
			Wrap(err)
	}

	return nil
}

func (s *cacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "Marshal").
			Wrap(err)
	}

	if err := s.redis.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		return domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "Set").
			Wrap(err)
	}

	s.logger.DebugContext(ctx, "cached key", "key", key, "ttl", ttl)
	return nil
}

func (s *cacheService) Delete(ctx context.Context, key string) error {
	if err := s.redis.Client.Del(ctx, key).Err(); err != nil {
		return domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "Delete").
			Wrap(err)
	}

	s.logger.DebugContext(ctx, "deleted cache key", "key", key)
	return nil
}

func (s *cacheService) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := s.redis.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "Exists").
			Wrap(err)
	}

	return exists > 0, nil
}

func (s *cacheService) RefreshTTL(ctx context.Context, key string, ttl time.Duration) error {
	if err := s.redis.Client.Expire(ctx, key, ttl).Err(); err != nil {
		s.logger.WarnContext(ctx, "failed to refresh cache TTL", "key", key, "error", err)
		return domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "RefreshTTL").
			Wrap(err)
	}

	return nil
}

func (s *cacheService) SetInvalidationFlag(ctx context.Context, key string, ttl time.Duration) error {
	if err := s.redis.Client.Set(ctx, key, "1", ttl).Err(); err != nil {
		return domain.ErrCacheError.
			WithContext("key", key).
			WithContext("operation", "SetInvalidationFlag").
			Wrap(err)
	}

	return nil
}
