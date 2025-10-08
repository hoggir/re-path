package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/app/http/response"
)

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
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
