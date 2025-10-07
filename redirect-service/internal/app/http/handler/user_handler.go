package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/app/http/response"
	"github.com/hoggir/re-path/redirect-service/internal/app/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserService.GetAll()
	if err != nil {
		response.InternalError(c, "Failed to fetch users", err.Error())
		return
	}

	response.OK(c, "Users fetched successfully", users)
}
