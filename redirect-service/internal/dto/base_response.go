package dto

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
)

type Response struct {
	Success   bool         `json:"success"`
	Message   string       `json:"message"`
	Data      interface{}  `json:"data,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

type ErrorDetail struct {
	Code     string                 `json:"code"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
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
		Success: false,
		Message: message,
		Error: &ErrorDetail{
			Code:    "ERROR",
			Message: message,
		},
		Timestamp: time.Now(),
	})
}

func HandleError(c *gin.Context, err error) {
	var appErr *domain.AppError

	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, Response{
			Success: false,
			Message: appErr.Message,
			Error: &ErrorDetail{
				Code:     appErr.Code,
				Message:  appErr.Message,
				Metadata: appErr.Metadata,
			},
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(500, Response{
		Success: false,
		Message: "An unexpected error occurred",
		Error: &ErrorDetail{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred",
		},
		Timestamp: time.Now(),
	})
}
