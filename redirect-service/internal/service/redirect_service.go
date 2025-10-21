package service

import (
	"context"
	"fmt"
	"log"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
)

type RedirectService interface {
	GetURL(ctx context.Context, shortCode string) (*domain.FindByShortCode, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
}

type redirectService struct {
	urlRepo      repository.URLRepository
	cacheService CacheService
	config       *config.Config
}

func NewRedirectService(
	urlRepo repository.URLRepository,
	cacheService CacheService,
	cfg *config.Config,
) RedirectService {
	return &redirectService{
		urlRepo:      urlRepo,
		cacheService: cacheService,
		config:       cfg,
	}
}

func (s *redirectService) GetURL(ctx context.Context, shortCode string) (*domain.FindByShortCode, error) {
	cacheKey := fmt.Sprintf("url:%s", shortCode)

	var url domain.FindByShortCode
	err := s.cacheService.Get(ctx, cacheKey, &url)
	if err == nil {
		dashboardCacheKey := fmt.Sprintf("dashboard:%d", url.UserID)
		log.Printf("⚡ Cache HIT for shortCode: %s", shortCode)
		s.cacheService.RefreshTTL(ctx, cacheKey, s.config.Redis.CacheTTL)
		s.cacheService.Delete(ctx, dashboardCacheKey)
		return &url, nil
	}

	urlData, err := s.urlRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if err := s.cacheService.Set(ctx, cacheKey, urlData, s.config.Redis.CacheTTL); err != nil {
		log.Printf("⚠️  Failed to cache shortCode %s: %v", shortCode, err)
	}

	dashboardCacheKey := fmt.Sprintf("dashboard:%d", url.UserID)
	s.cacheService.Delete(ctx, dashboardCacheKey)

	return urlData, nil
}

func (s *redirectService) IncrementClickCount(ctx context.Context, shortCode string) error {
	return s.urlRepo.IncrementClickCount(ctx, shortCode)
}
