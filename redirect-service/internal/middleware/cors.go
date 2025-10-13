package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/config"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origins := strings.Split(cfg.CORS.AllowOrigins, ",")
		methods := strings.Split(cfg.CORS.AllowMethods, ",")
		headers := strings.Split(cfg.CORS.AllowHeaders, ",")

		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, allowedOrigin := range origins {
			allowedOrigin = strings.TrimSpace(allowedOrigin)
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(origins) > 0 && origins[0] == "*" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		c.Writer.Header().Set("Access-Control-Expose-Headers", "Location")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
