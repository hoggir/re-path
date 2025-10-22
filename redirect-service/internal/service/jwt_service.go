package service

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
)

type JWTService interface {
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jwtService struct {
	config *config.Config
}

type JWTClaims struct {
	Sub   interface{} `json:"sub"`
	Email string      `json:"email"`
	Role  string      `json:"role"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) GetUserID() string {
	switch v := c.Sub.(type) {
	case float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (c *JWTClaims) GetUserIDAsInt() int {
	switch v := c.Sub.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		// If it's a string, try to parse it, otherwise return 0
		var id int
		fmt.Sscanf(v, "%d", &id)
		return id
	default:
		return 0
	}
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		config: cfg,
	}
}

func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidSigningKey
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired.Wrap(err)
		}
		return nil, domain.ErrInvalidToken.Wrap(err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}
