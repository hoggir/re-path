package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/app/http/handler"
)

func NewRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	api.GET("/users", userHandler.GetAllUsers)

	return r
}
