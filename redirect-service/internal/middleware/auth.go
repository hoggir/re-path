package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/dto"
	"github.com/hoggir/re-path/redirect-service/internal/service"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	UserEmailKey        = "user_email"
	UserRoleKey         = "user_role"
)

func JWTAuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			dto.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			dto.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format. Expected Bearer token", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			dto.ErrorResponse(c, http.StatusUnauthorized, "Token is required", nil)
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			switch err {
			case service.ErrExpiredToken:
				dto.ErrorResponse(c, http.StatusUnauthorized, "Token has expired", err.Error())
			case service.ErrInvalidToken:
				dto.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", err.Error())
			default:
				dto.ErrorResponse(c, http.StatusUnauthorized, "Token validation failed", err.Error())
			}
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.GetUserIDAsInt())
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}
