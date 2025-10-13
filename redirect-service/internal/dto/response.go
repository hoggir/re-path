package dto

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	})
}

type RedirectResponse struct {
	ShortCode   string `json:"shortCode"`
	OriginalURL string `json:"originalUrl"`
	ClickCount  int    `json:"clickCount"`
}
