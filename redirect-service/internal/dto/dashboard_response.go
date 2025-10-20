package dto

type DashboardResponse struct {
	TotalLink      int `json:"total_link" example:"666"`
	TotalClick     int `json:"total_click" example:"333"`
	UniqueVisitors int `json:"unique_visitors" example:"123"`
}
