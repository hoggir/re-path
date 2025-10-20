package domain

import (
	"errors"
	"fmt"
)

// DashboardRequest represents RPC request for dashboard data
type DashboardRequest struct {
	UserID int `json:"user_id" validate:"required,gt=0"`
}

// Validate checks if the request is valid
func (r *DashboardRequest) Validate() error {
	if r.UserID <= 0 {
		return errors.New("user_id must be greater than 0")
	}
	return nil
}

// RecentClick represents a recent click event in dashboard
type RecentClick struct {
	ShortCode     string  `json:"short_code"`
	ClickedAt     string  `json:"clicked_at"`
	IPAddressHash *string `json:"ip_address_hash,omitempty"`
	UserAgent     *string `json:"user_agent,omitempty"`
	CountryCode   *string `json:"country_code,omitempty"`
	City          *string `json:"city,omitempty"`
	DeviceType    *string `json:"device_type,omitempty"`
	BrowserName   *string `json:"browser_name,omitempty"`
	IsBot         bool    `json:"is_bot"`
}

// TopLink represents top link by click count
type TopLink struct {
	ShortCode string `json:"short_code"`
	Clicks    int    `json:"clicks" validate:"gte=0"`
}

// DashboardResponse represents RPC response from analytic-service
type DashboardResponse struct {
	UserID       int           `json:"user_id"`
	TotalClicks  int           `json:"total_clicks" validate:"gte=0"`
	TotalLinks   int           `json:"total_links" validate:"gte=0"`
	RecentClicks []RecentClick `json:"recent_clicks"`
	TopLinks     []TopLink     `json:"top_links"`
	Status       string        `json:"status" validate:"required,oneof=success error limited"`
	Message      *string       `json:"message,omitempty"`
}

// Validate checks if the response is valid
func (r *DashboardResponse) Validate() error {
	if r.UserID <= 0 {
		return fmt.Errorf("invalid user_id: %d", r.UserID)
	}
	if r.TotalClicks < 0 {
		return fmt.Errorf("invalid total_clicks: %d", r.TotalClicks)
	}
	if r.TotalLinks < 0 {
		return fmt.Errorf("invalid total_links: %d", r.TotalLinks)
	}
	if r.Status != "success" && r.Status != "error" && r.Status != "limited" {
		return fmt.Errorf("invalid status: %s", r.Status)
	}
	if len(r.RecentClicks) > 10 {
		return fmt.Errorf("recent_clicks exceeds limit: %d", len(r.RecentClicks))
	}
	if len(r.TopLinks) > 5 {
		return fmt.Errorf("top_links exceeds limit: %d", len(r.TopLinks))
	}
	return nil
}

// IsSuccess checks if response status is success
func (r *DashboardResponse) IsSuccess() bool {
	return r.Status == "success"
}

// IsError checks if response status is error
func (r *DashboardResponse) IsError() bool {
	return r.Status == "error"
}

// IsLimited checks if response status is limited
func (r *DashboardResponse) IsLimited() bool {
	return r.Status == "limited"
}

// GetMessage safely retrieves the message
func (r *DashboardResponse) GetMessage() string {
	if r.Message != nil {
		return *r.Message
	}
	return ""
}
