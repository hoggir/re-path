//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/app/http"
	"github.com/hoggir/re-path/redirect-service/internal/app/http/handler"
	"github.com/hoggir/re-path/redirect-service/internal/app/repository"
	"github.com/hoggir/re-path/redirect-service/internal/app/service"
	"github.com/google/wire"
)

func InitializeApp() *gin.Engine {
	wire.Build(
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,
		http.NewRouter,
	)
	return nil
}
