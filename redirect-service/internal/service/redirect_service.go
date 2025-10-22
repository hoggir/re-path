package service

import (
	"context"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
)

type RedirectService interface {
	GetURL(ctx context.Context, shortCode string) (*domain.FindByShortCode, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
}

type redirectService struct {
	urlRepo      repository.URLRepository
	cacheService CacheService
	cacheKeys    *CacheKeyGenerator
	config       *config.Config
	logger       logger.Logger
}

func NewRedirectService(
	urlRepo repository.URLRepository,
	cacheService CacheService,
	cacheKeys *CacheKeyGenerator,
	cfg *config.Config,
	log logger.Logger,
) RedirectService {
	return &redirectService{
		urlRepo:      urlRepo,
		cacheService: cacheService,
		cacheKeys:    cacheKeys,
		config:       cfg,
		logger:       log,
	}
}

func (s *redirectService) GetURL(ctx context.Context, shortCode string) (*domain.FindByShortCode, error) {
	cacheKey := s.cacheKeys.URL(shortCode)

	var url domain.FindByShortCode
	err := s.cacheService.Get(ctx, cacheKey, &url)
	if err == nil {
		dashboardInvalidFlag := s.cacheKeys.DashboardInvalidationFlag(url.UserID)
		s.logger.DebugContext(ctx, "cache hit for shortCode", "shortCode", shortCode)
		s.cacheService.RefreshTTL(ctx, cacheKey, s.config.Redis.CacheTTL)
		s.cacheService.SetInvalidationFlag(ctx, dashboardInvalidFlag, s.config.Redis.InvalidationFlagTTL)
		return &url, nil
	}

	urlData, err := s.urlRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if err := s.cacheService.Set(ctx, cacheKey, urlData, s.config.Redis.CacheTTL); err != nil {
		s.logger.WarnContext(ctx, "failed to cache shortCode", "shortCode", shortCode, "error", err)
	}

	dashboardInvalidFlag := s.cacheKeys.DashboardInvalidationFlag(urlData.UserID)
	s.cacheService.SetInvalidationFlag(ctx, dashboardInvalidFlag, s.config.Redis.InvalidationFlagTTL)

	return urlData, nil
}

func (s *redirectService) IncrementClickCount(ctx context.Context, shortCode string) error {
	return s.urlRepo.IncrementClickCount(ctx, shortCode)
}
