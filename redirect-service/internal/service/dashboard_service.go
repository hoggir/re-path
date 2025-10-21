package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, userId int) (*domain.DashboardResponse, error)
}

type dashboardService struct {
	rpcService   RabbitMQRPCService
	cacheService CacheService
	config       *config.Config
}

func NewDashboardService(rpcService RabbitMQRPCService, cacheService CacheService, cfg *config.Config) DashboardService {
	return &dashboardService{
		rpcService:   rpcService,
		cacheService: cacheService,
		config:       cfg,
	}
}

func (s *dashboardService) GetDashboard(ctx context.Context, userId int) (*domain.DashboardResponse, error) {
	cacheKey := fmt.Sprintf("dashboard:%d", userId)

	var cachedResponse domain.DashboardResponse
	if err := s.cacheService.Get(ctx, cacheKey, &cachedResponse); err == nil {
		log.Printf("⚡ Cache HIT for dashboard: %d", userId)
		s.cacheService.RefreshTTL(ctx, cacheKey, s.config.Redis.CacheTTL)
		return &cachedResponse, nil
	}

	request := domain.DashboardRequest{
		UserID: userId,
	}

	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	response, err := s.rpcService.Call(
		ctx,
		s.config.RabbitMQ.Queues.DashboardRequest,
		request,
		5*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard data: %w", err)
	}

	var result domain.DashboardResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to parse dashboard response: %w", err)
	}

	if result.IsError() {
		return nil, fmt.Errorf("dashboard service error: %s", result.GetMessage())
	}

	if err := s.cacheService.Set(ctx, cacheKey, result, s.config.Redis.CacheTTL); err != nil {
		log.Printf("⚠️  Failed to cache dashboard: %d: %v", userId, err)
	}

	return &result, nil
}
