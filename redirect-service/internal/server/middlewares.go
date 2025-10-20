package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/middleware"
	"github.com/hoggir/re-path/redirect-service/internal/service"
)

// Middlewares groups all middleware functions
type Middlewares struct {
	CORS gin.HandlerFunc
	Auth gin.HandlerFunc
}

// NewMiddlewares creates a new Middlewares instance
func NewMiddlewares(cfg *config.Config, jwtService service.JWTService) *Middlewares {
	return &Middlewares{
		CORS: middleware.CORSMiddleware(cfg),
		Auth: middleware.JWTAuthMiddleware(jwtService),
	}
}
