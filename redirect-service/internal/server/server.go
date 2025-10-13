package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/handler"
	"github.com/hoggir/re-path/redirect-service/internal/middleware"
)

type Server struct {
	Config          *config.Config
	Router          *gin.Engine
	RedirectHandler *handler.RedirectHandler
	HealthHandler   *handler.HealthHandler
	MongoDB         *database.MongoDB
	Redis           *database.Redis
}

func New(
	cfg *config.Config,
	redirectHandler *handler.RedirectHandler,
	healthHandler *handler.HealthHandler,
	mongoDB *database.MongoDB,
	redis *database.Redis,
) *Server {
	gin.SetMode(cfg.Server.GinMode)

	srv := &Server{
		Config:          cfg,
		RedirectHandler: redirectHandler,
		HealthHandler:   healthHandler,
		MongoDB:         mongoDB,
		Redis:           redis,
	}

	srv.setupRouter()

	return srv
}

func (s *Server) setupRouter() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(s.Config))

	s.registerRoutes(r)

	s.Router = r
}

func (s *Server) GetRouter() *gin.Engine {
	return s.Router
}
