//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/hoggir/re-path/redirect-service/internal/app/http"
	"github.com/hoggir/re-path/redirect-service/internal/app/http/handler"
	usersPkg "github.com/hoggir/re-path/redirect-service/internal/users"
)

func InitializeApp() *gin.Engine {
	wire.Build(
		usersPkg.NewUserRepository,
		usersPkg.NewUserService,
		handler.NewUserHandler,
		http.NewRouter,
	)
	return nil
}
