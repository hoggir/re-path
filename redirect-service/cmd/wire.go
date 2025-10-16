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

		repository.NewURLRepository,
		repository.NewClickEventRepository,

		service.NewCacheService,
		service.NewGeoIPService,
		service.NewRedirectService,
		service.NewClickEventService,

		handler.NewRedirectHandler,
		handler.NewHealthHandler,

		server.New,
	)
	return nil, nil
}
