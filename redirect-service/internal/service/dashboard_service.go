package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, userId int) (*domain.DashboardResponse, error)
}

type dashboardService struct {
	rpcService   RabbitMQRPCService
	cacheService CacheService
	cacheKeys    *CacheKeyGenerator
	config       *config.Config
	logger       logger.Logger
}

func NewDashboardService(rpcService RabbitMQRPCService, cacheService CacheService, cacheKeys *CacheKeyGenerator, cfg *config.Config, log logger.Logger) DashboardService {
	return &dashboardService{
		rpcService:   rpcService,
		cacheService: cacheService,
		cacheKeys:    cacheKeys,
		config:       cfg,
		logger:       log,
	}
}

func (s *dashboardService) GetDashboard(ctx context.Context, userId int) (*domain.DashboardResponse, error) {
	cacheKey := s.cacheKeys.Dashboard(userId)
	invalidFlagKey := s.cacheKeys.DashboardInvalidationFlag(userId)

	invalidFlagExists, err := s.cacheService.Exists(ctx, invalidFlagKey)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to check invalidation flag", "userId", userId, "error", err)
	}

	if !invalidFlagExists {
		var cachedResponse domain.DashboardResponse
		if err := s.cacheService.Get(ctx, cacheKey, &cachedResponse); err == nil {
			s.logger.DebugContext(ctx, "cache hit for dashboard", "userId", userId)
			s.cacheService.RefreshTTL(ctx, cacheKey, s.config.Redis.CacheTTL)
			return &cachedResponse, nil
		}
	} else {
		s.logger.DebugContext(ctx, "dashboard invalidation flag found, refreshing from RPC", "userId", userId)
		s.cacheService.Delete(ctx, invalidFlagKey)
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
		s.config.RabbitMQ.RPCTimeout,
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
		s.logger.WarnContext(ctx, "failed to cache dashboard", "userId", userId, "error", err)
	}

	return &result, nil
}
