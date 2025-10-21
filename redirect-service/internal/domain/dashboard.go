package domain

import (
	"errors"
)

type DashboardRequest struct {
	UserID int `json:"user_id" validate:"required,gt=0"`
}

func (r *DashboardRequest) Validate() error {
	if r.UserID <= 0 {
		return errors.New("user_id must be greater than 0")
	}
	return nil
}

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

type TopLink struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
	Clicks      int    `json:"clicks"`
	Status      bool   `json:"status"`
}

type StatLink struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

type DashboardResponse struct {
	UserID       int        `json:"user_id"`
	TotalClicks  int        `json:"total_clicks" validate:"gte=0"`
	TotalLinks   int        `json:"total_links" validate:"gte=0"`
	UniqVisitors int        `json:"uniq_visitors"`
	TopLinks     []TopLink  `json:"top_links"`
	StatLinks    []StatLink `json:"stat_links"`
	Status       string     `json:"status" validate:"required,oneof=success error limited"`
	Message      *string    `json:"message,omitempty"`
}

func (r *DashboardResponse) IsSuccess() bool {
	return r.Status == "success"
}

func (r *DashboardResponse) IsError() bool {
	return r.Status == "error"
}

func (r *DashboardResponse) IsLimited() bool {
	return r.Status == "limited"
}

func (r *DashboardResponse) GetMessage() string {
	if r.Message != nil {
		return *r.Message
	}
	return ""
}
