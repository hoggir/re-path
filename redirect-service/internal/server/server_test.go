package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/handler"
	"github.com/hoggir/re-path/redirect-service/internal/server"
)

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			GinMode: "test",
		},
	}

	healthHandler := handler.NewHealthHandler()

	srv := server.New(cfg, nil, healthHandler, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	srv.GetRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
