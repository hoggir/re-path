package dto

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status" example:"UP"`
	Service string `json:"service" example:"redirect-service"`
	Version string `json:"version,omitempty" example:"1.0.0"`
}
