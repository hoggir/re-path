package dto

import "github.com/hoggir/re-path/redirect-service/internal/domain"

type DashboardResponse struct {
	TotalLink    int               `json:"total_link" example:"666"`
	TotalClick   int               `json:"total_click" example:"333"`
	UniqVisitors int               `json:"uniq_visitors" example:"2"`
	TopLinks     []domain.TopLink  `json:"top_links"`
	StatLinks    []domain.StatLink `json:"stat_links"`
}
