package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
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
			dto.HandleError(c, domain.ErrUnauthorized.WithMessage("Authorization header is required"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			dto.HandleError(c, domain.ErrUnauthorized.WithMessage("Invalid authorization format. Expected Bearer token"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			dto.HandleError(c, domain.ErrUnauthorized.WithMessage("Token is required"))
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			dto.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.GetUserIDAsInt())
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}
