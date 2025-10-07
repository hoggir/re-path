package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ---- SUCCESS WRAPPER ----
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, SuccessResponse{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ---- ERROR WRAPPER ----
func Error(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: message,
		Errors:  err,
	})
}

// ---- Shortcut Helpers ----
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

func BadRequest(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusBadRequest, message, err)
}

func NotFound(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusNotFound, message, err)
}

func InternalError(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusInternalServerError, message, err)
}
