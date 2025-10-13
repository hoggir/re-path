package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
	"github.com/redis/go-redis/v9"
)

type RedirectService interface {
	GetURL(ctx context.Context, shortCode string) (*domain.URL, error)
}

type redirectService struct {
	urlRepo repository.URLRepository
	redis   *database.Redis
	config  *config.Config
}

func NewRedirectService(
	urlRepo repository.URLRepository,
	redis *database.Redis,
	cfg *config.Config,
) RedirectService {
	return &redirectService{
		urlRepo: urlRepo,
		redis:   redis,
		config:  cfg,
	}
}

func (s *redirectService) GetURL(ctx context.Context, shortCode string) (*domain.URL, error) {
	cacheKey := fmt.Sprintf("url:%s", shortCode)

	cachedData, err := s.redis.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var url domain.URL
		if err := json.Unmarshal([]byte(cachedData), &url); err == nil {
			log.Printf("⚡ Cache HIT for shortCode: %s", shortCode)

			return &url, nil
		}
	} else if err != redis.Nil {
		log.Printf("⚠️  Redis error for shortCode %s: %v", shortCode, err)
	}

	log.Printf("💾 Cache MISS for shortCode: %s, querying database...", shortCode)
	url, err := s.urlRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		// Langsung return error dari repository (no redundant wrapping)
		return nil, err
	}

	urlJSON, err := json.Marshal(url)
	if err == nil {
		if err := s.redis.Client.Set(ctx, cacheKey, urlJSON, s.config.Redis.CacheTTL).Err(); err != nil {
			log.Printf("⚠️  Failed to cache shortCode %s: %v", shortCode, err)
		} else {
			log.Printf("✅ Cached shortCode: %s for %v", shortCode, s.config.Redis.CacheTTL)
		}
	}

	return url, nil
}
