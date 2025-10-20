package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hoggir/re-path/redirect-service/internal/domain"
)

// ExampleDashboardRPC demonstrates how to use typed RPC calls
func ExampleDashboardRPC() {
	// 1. Create and validate request
	request := domain.DashboardRequest{
		UserID: 123,
	}

	if err := request.Validate(); err != nil {
		log.Fatalf("Invalid request: %v", err)
	}

	// 2. Marshal to JSON (this would be sent via RabbitMQ)
	requestJSON, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	fmt.Printf("Request JSON: %s\n", string(requestJSON))

	// 3. Simulate RPC call (in real code, use RabbitMQRPCService)
	// response := rpcService.Call(ctx, "dashboard_request", request, timeout)

	// 4. Parse response from Python service
	pythonResponse := `{
		"user_id": 123,
		"total_clicks": 1542,
		"total_links": 45,
		"recent_clicks": [
			{
				"short_code": "my-link",
				"clicked_at": "2025-10-20T10:00:00Z",
				"ip_address_hash": "abc123",
				"country_code": "ID",
				"city": "Jakarta",
				"is_bot": false
			}
		],
		"top_links": [
			{"short_code": "popular", "clicks": 350}
		],
		"status": "success"
	}`

	var response domain.DashboardResponse
	if err := json.Unmarshal([]byte(pythonResponse), &response); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	// 5. Validate response
	if err := response.Validate(); err != nil {
		log.Fatalf("Invalid response: %v", err)
	}

	// 6. Handle response based on status
	switch {
	case response.IsSuccess():
		fmt.Printf("✅ Success! Total clicks: %d, Total links: %d\n",
			response.TotalClicks, response.TotalLinks)

		// Access typed data
		for _, click := range response.RecentClicks {
			fmt.Printf("  - %s clicked at %s\n", click.ShortCode, click.ClickedAt)
		}

		for _, link := range response.TopLinks {
			fmt.Printf("  - Top link: %s (%d clicks)\n", link.ShortCode, link.Clicks)
		}

	case response.IsError():
		log.Printf("❌ Error from analytic service: %s", response.GetMessage())

	case response.IsLimited():
		log.Printf("⚠️  Limited data: %s", response.GetMessage())
		// Still usable but might be incomplete
		fmt.Printf("Partial data - Clicks: %d, Links: %d\n",
			response.TotalClicks, response.TotalLinks)
	}
}

// ExampleWithValidator demonstrates using the ContractValidator utility
func ExampleWithValidator() {
	ctx := context.Background()
	_ = ctx // Use in actual RPC call

	// Using validator for automatic validation
	requestJSON := []byte(`{"user_id": 123}`)

	// Validate and parse request
	request, err := domain.Validator.ValidateDashboardRequestJSON(requestJSON)
	if err != nil {
		log.Fatalf("Invalid request: %v", err)
	}

	fmt.Printf("Validated request: UserID=%d\n", request.UserID)

	// Simulate response from Python
	responseJSON := []byte(`{
		"user_id": 123,
		"total_clicks": 100,
		"total_links": 10,
		"recent_clicks": [],
		"top_links": [],
		"status": "success"
	}`)

	// Validate and parse response
	response, err := domain.Validator.ValidateDashboardResponseJSON(responseJSON)
	if err != nil {
		log.Fatalf("Invalid response: %v", err)
	}

	fmt.Printf("Validated response: Status=%s, Clicks=%d\n",
		response.Status, response.TotalClicks)
}

// ExampleErrorHandling demonstrates proper error handling
func ExampleErrorHandling() {
	// Example 1: Invalid request
	invalidRequest := domain.DashboardRequest{UserID: 0}
	if err := invalidRequest.Validate(); err != nil {
		fmt.Printf("Caught invalid request: %v\n", err)
		// Handle error appropriately
	}

	// Example 2: Error response from service
	errorResponseJSON := `{
		"user_id": 123,
		"status": "error",
		"message": "Database connection failed",
		"total_clicks": 0,
		"total_links": 0,
		"recent_clicks": [],
		"top_links": []
	}`

	var response domain.DashboardResponse
	if err := json.Unmarshal([]byte(errorResponseJSON), &response); err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	// Validate structure
	if err := response.Validate(); err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	// Check status
	if response.IsError() {
		// Handle error status
		fmt.Printf("Service returned error: %s\n", response.GetMessage())
		// Maybe retry, or return error to user
	}

	// Example 3: Invalid response from service
	invalidResponseJSON := `{
		"user_id": -1,
		"status": "success",
		"total_clicks": -100,
		"total_links": 0,
		"recent_clicks": [],
		"top_links": []
	}`

	var invalidResponse domain.DashboardResponse
	json.Unmarshal([]byte(invalidResponseJSON), &invalidResponse)

	if err := invalidResponse.Validate(); err != nil {
		// Validation catches the invalid data
		fmt.Printf("Caught invalid response from service: %v\n", err)
		// Log error, alert monitoring, etc.
	}
}

// ExampleServiceIntegration shows integration with actual service layer
func ExampleServiceIntegration(dashboardService interface {
	GetDashboard(ctx context.Context, userId int) (*domain.DashboardResponse, error)
}) {
	ctx := context.Background()
	userID := 123

	// Service call with automatic validation
	response, err := dashboardService.GetDashboard(ctx, userID)
	if err != nil {
		log.Printf("Dashboard service error: %v", err)
		return
	}

	// Response is already validated by service layer
	// Safe to use directly
	fmt.Printf("Dashboard loaded: %d clicks, %d links\n",
		response.TotalClicks, response.TotalLinks)

	// Type-safe access to nested data
	if len(response.RecentClicks) > 0 {
		latestClick := response.RecentClicks[0]
		fmt.Printf("Latest click: %s at %s\n",
			latestClick.ShortCode, latestClick.ClickedAt)

		// Safe nil pointer access
		if latestClick.CountryCode != nil {
			fmt.Printf("Country: %s\n", *latestClick.CountryCode)
		}
	}

	// Helper methods
	if response.IsLimited() {
		log.Printf("Warning: %s", response.GetMessage())
	}
}
