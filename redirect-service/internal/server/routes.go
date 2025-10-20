package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/hoggir/re-path/redirect-service/docs"
)

func (s *Server) registerRoutes(r *gin.Engine) {
	s.registerPublicRoutes(r)
	s.registerAPIRoutes(r)
	s.registerRedirectRoutes(r)
}

func (s *Server) registerPublicRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", s.Handlers.Health.Health)
}

func (s *Server) registerAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/info/:shortUrl", s.Handlers.Redirect.GetURLInfo)
		s.registerProtectedAPIRoutes(api)
	}
}

func (s *Server) registerProtectedAPIRoutes(rg *gin.RouterGroup) {
	protected := rg.Group("")
	protected.Use(s.Middlewares.Auth)
	{
		protected.GET("/dashboard", s.Handlers.Dashboard.GetDashboardByShortUrl)
	}
}

func (s *Server) registerRedirectRoutes(r *gin.Engine) {
	r.GET("/r/:shortUrl", s.Handlers.Redirect.Redirect)
}
