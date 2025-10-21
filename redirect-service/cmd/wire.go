//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/handler"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
	"github.com/hoggir/re-path/redirect-service/internal/server"
	"github.com/hoggir/re-path/redirect-service/internal/service"
)

func InitializeApp() (*server.Server, error) {
	wire.Build(
		config.Load,

		database.NewMongoDB,
		database.NewRedis,
		database.NewRabbitMQ,

		repository.NewURLRepository,
		repository.NewClickEventRepository,

		service.NewCacheService,
		service.NewGeoIPService,
		service.NewRabbitMQRPCService,
		service.NewRedirectService,
		service.NewClickEventService,
		service.NewDashboardService,
		service.NewJWTService,

		handler.NewRedirectHandler,
		handler.NewHealthHandler,
		handler.NewDashboardHandler,

		server.NewHandlers,
		server.NewMiddlewares,
		server.New,
	)
	return nil, nil
}
