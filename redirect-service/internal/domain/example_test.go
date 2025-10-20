package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestParseRequestExample demonstrates parsing request example JSON
func TestParseRequestExample(t *testing.T) {
	// Load example JSON from analytic-service
	examplePath := filepath.Join("..", "..", "..", "analytic-service", "schemas", "dashboard_request.example.json")

	// Skip if file doesn't exist (e.g., in CI without analytic-service)
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("Example file not found - skipping")
	}

	data, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to read example: %v", err)
	}

	// Parse into Go struct
	var request DashboardRequest
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Validate
	if err := request.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Verify structure
	if request.UserID <= 0 {
		t.Errorf("UserID should be > 0, got %d", request.UserID)
	}

	t.Logf("✅ Request example parsed successfully")
	t.Logf("   UserID: %d", request.UserID)
}

// TestParseResponseExample demonstrates parsing response example JSON
func TestParseResponseExample(t *testing.T) {
	examplePath := filepath.Join("..", "..", "..", "analytic-service", "schemas", "dashboard_response.example.json")

	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("Example file not found - skipping")
	}

	data, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to read example: %v", err)
	}

	var response DashboardResponse
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Validate structure
	if err := response.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Verify expected values
	if response.UserID <= 0 {
		t.Errorf("UserID should be > 0, got %d", response.UserID)
	}

	if response.TotalClicks < 0 {
		t.Errorf("TotalClicks should be >= 0, got %d", response.TotalClicks)
	}

	if response.TotalLinks < 0 {
		t.Errorf("TotalLinks should be >= 0, got %d", response.TotalLinks)
	}

	if !response.IsSuccess() {
		t.Errorf("Example should have success status, got %s", response.Status)
	}

	t.Logf("✅ Response example parsed successfully")
	t.Logf("   UserID: %d", response.UserID)
	t.Logf("   TotalClicks: %d", response.TotalClicks)
	t.Logf("   TotalLinks: %d", response.TotalLinks)
	t.Logf("   RecentClicks count: %d", len(response.RecentClicks))
	t.Logf("   TopLinks count: %d", len(response.TopLinks))
	t.Logf("   Status: %s", response.Status)
}

// TestParseErrorExample demonstrates parsing error example JSON
func TestParseErrorExample(t *testing.T) {
	examplePath := filepath.Join("..", "..", "..", "analytic-service", "schemas", "dashboard_error.example.json")

	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("Example file not found - skipping")
	}

	data, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to read example: %v", err)
	}

	var response DashboardResponse
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Validate structure
	if err := response.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Verify it's an error response
	if !response.IsError() {
		t.Errorf("Example should have error status, got %s", response.Status)
	}

	// Verify message exists
	if response.GetMessage() == "" {
		t.Errorf("Error response should have message")
	}

	t.Logf("✅ Error example parsed successfully")
	t.Logf("   Status: %s", response.Status)
	t.Logf("   Message: %s", response.GetMessage())
	t.Logf("   TotalClicks: %d", response.TotalClicks)
	t.Logf("   TotalLinks: %d", response.TotalLinks)
}

// TestExampleJSONRoundTrip demonstrates full round-trip
func TestExampleJSONRoundTrip(t *testing.T) {
	// Create example response
	message := "Test message"
	response := DashboardResponse{
		UserID:      123,
		TotalClicks: 1000,
		TotalLinks:  50,
		RecentClicks: []RecentClick{
			{
				ShortCode: "test-link",
				ClickedAt: "2025-10-20T10:00:00Z",
				IsBot:     false,
			},
		},
		TopLinks: []TopLink{
			{ShortCode: "popular", Clicks: 500},
		},
		Status:  "success",
		Message: &message,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	t.Logf("JSON output:\n%s", string(jsonData))

	// Unmarshal back
	var decoded DashboardResponse
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify
	if decoded.UserID != response.UserID {
		t.Errorf("UserID mismatch: got %d, want %d", decoded.UserID, response.UserID)
	}

	if decoded.TotalClicks != response.TotalClicks {
		t.Errorf("TotalClicks mismatch: got %d, want %d", decoded.TotalClicks, response.TotalClicks)
	}

	if decoded.Status != response.Status {
		t.Errorf("Status mismatch: got %s, want %s", decoded.Status, response.Status)
	}

	t.Logf("✅ Round-trip successful")
}

// ExampleDashboardRequest_Unmarshal shows how to use request example
func ExampleDashboardRequest_Unmarshal() {
	// In real code, you would load from file:
	// data, _ := os.ReadFile("schemas/dashboard_request.example.json")

	jsonData := []byte(`{"user_id": 123}`)

	var request DashboardRequest
	if err := json.Unmarshal(jsonData, &request); err != nil {
		panic(err)
	}

	if err := request.Validate(); err != nil {
		panic(err)
	}

	// Use request
	_ = request.UserID
	// Output:
}

// ExampleDashboardResponse_Unmarshal shows how to use response example
func ExampleDashboardResponse_Unmarshal() {
	// In real code, you would load from file:
	// data, _ := os.ReadFile("schemas/dashboard_response.example.json")

	jsonData := []byte(`{
		"user_id": 123,
		"total_clicks": 1000,
		"total_links": 50,
		"recent_clicks": [],
		"top_links": [],
		"status": "success"
	}`)

	var response DashboardResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		panic(err)
	}

	if err := response.Validate(); err != nil {
		panic(err)
	}

	// Check status
	if response.IsSuccess() {
		// Process success response
		_ = response.TotalClicks
	}
	// Output:
}
