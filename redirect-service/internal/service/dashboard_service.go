package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, userId int) (*domain.DashboardResponse, error)
}

type dashboardService struct {
	rpcService RabbitMQRPCService
	config     *config.Config
}

func NewDashboardService(rpcService RabbitMQRPCService, cfg *config.Config) DashboardService {
	return &dashboardService{
		rpcService: rpcService,
		config:     cfg,
	}
}

func (s *dashboardService) GetDashboard(ctx context.Context, userId int) (*domain.DashboardResponse, error) {
	// Create and validate request
	request := domain.DashboardRequest{
		UserID: userId,
	}

	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Make RPC call
	response, err := s.rpcService.Call(
		ctx,
		s.config.RabbitMQ.Queues.DashboardRequest,
		request,
		5*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard data: %w", err)
	}

	// Parse response into typed struct
	var result domain.DashboardResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to parse dashboard response: %w", err)
	}

	// Validate response
	if err := result.Validate(); err != nil {
		return nil, fmt.Errorf("invalid response from analytic service: %w", err)
	}

	// Check response status
	if result.IsError() {
		return nil, fmt.Errorf("dashboard service error: %s", result.GetMessage())
	}

	return &result, nil
}
