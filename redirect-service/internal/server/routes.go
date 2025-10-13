package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/hoggir/re-path/redirect-service/docs"
)

func (s *Server) registerRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", s.HealthHandler.Health)

	api := r.Group("/api")
	{
		api.GET("/info/:shortCode", s.RedirectHandler.GetURLInfo)
	}

	r.GET("/:shortCode", s.RedirectHandler.Redirect)
}
