package dto

import "time"

type RedirectResponse struct {
	OriginalURL string `json:"originalUrl" example:"https://example.com"`
}

type URLInfoResponse struct {
	ShortCode   string     `json:"shortCode" example:"abc123"`
	OriginalURL string     `json:"originalUrl" example:"https://example.com"`
	CustomAlias string     `json:"customAlias,omitempty" example:"my-link"`
	ClickCount  int        `json:"clickCount" example:"42"`
	IsActive    bool       `json:"isActive" example:"true"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty" example:"2024-12-31T23:59:59Z"`
	CreatedAt   time.Time  `json:"createdAt" example:"2024-01-01T00:00:00Z"`
}
