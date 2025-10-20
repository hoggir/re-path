package server

import "github.com/hoggir/re-path/redirect-service/internal/handler"

// Handlers groups all HTTP handlers
type Handlers struct {
	Redirect  *handler.RedirectHandler
	Health    *handler.HealthHandler
	Dashboard *handler.DashboardHandler
}

// NewHandlers creates a new Handlers instance
func NewHandlers(
	redirectHandler *handler.RedirectHandler,
	healthHandler *handler.HealthHandler,
	dashboardHandler *handler.DashboardHandler,
) *Handlers {
	return &Handlers{
		Redirect:  redirectHandler,
		Health:    healthHandler,
		Dashboard: dashboardHandler,
	}
}
