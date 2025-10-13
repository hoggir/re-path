package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/dto"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health check endpoint
// @Summary Health check
// @Tags Health
// @Success 200 {object} dto.Response
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	dto.SuccessResponse(c, http.StatusOK, "Service is healthy", gin.H{
		"status": "UP",
		"service": "redirect-service",
	})
}
