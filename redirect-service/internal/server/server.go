package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/database"
)

type Server struct {
	Config      *config.Config
	Router      *gin.Engine
	Handlers    *Handlers
	Middlewares *Middlewares
	MongoDB     *database.MongoDB
	Redis       *database.Redis
}

func New(
	cfg *config.Config,
	handlers *Handlers,
	middlewares *Middlewares,
	mongoDB *database.MongoDB,
	redis *database.Redis,
) *Server {
	gin.SetMode(cfg.Server.GinMode)

	srv := &Server{
		Config:      cfg,
		Handlers:    handlers,
		Middlewares: middlewares,
		MongoDB:     mongoDB,
		Redis:       redis,
	}

	srv.setupRouter()

	return srv
}

func (s *Server) setupRouter() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(s.Middlewares.CORS)

	s.registerRoutes(r)

	s.Router = r
}

func (s *Server) GetRouter() *gin.Engine {
	return s.Router
}
