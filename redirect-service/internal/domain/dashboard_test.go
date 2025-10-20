package domain

import (
	"encoding/json"
	"testing"
)

func TestDashboardRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request DashboardRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			request: DashboardRequest{UserID: 1},
			wantErr: false,
		},
		{
			name:    "valid request with large ID",
			request: DashboardRequest{UserID: 999999},
			wantErr: false,
		},
		{
			name:    "invalid - zero user_id",
			request: DashboardRequest{UserID: 0},
			wantErr: true,
		},
		{
			name:    "invalid - negative user_id",
			request: DashboardRequest{UserID: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DashboardRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDashboardRequest_JSON(t *testing.T) {
	request := DashboardRequest{UserID: 123}

	// Marshal to JSON
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Check JSON format
	expected := `{"user_id":123}`
	if string(data) != expected {
		t.Errorf("JSON mismatch: got %s, want %s", string(data), expected)
	}

	// Unmarshal back
	var decoded DashboardRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if decoded.UserID != request.UserID {
		t.Errorf("Decoded user_id = %d, want %d", decoded.UserID, request.UserID)
	}
}

func TestDashboardResponse_Validate(t *testing.T) {
	validMessage := "test message"

	tests := []struct {
		name     string
		response DashboardResponse
		wantErr  bool
	}{
		{
			name: "valid success response",
			response: DashboardResponse{
				UserID:       1,
				TotalClicks:  100,
				TotalLinks:   10,
				RecentClicks: []RecentClick{},
				TopLinks:     []TopLink{},
				Status:       "success",
			},
			wantErr: false,
		},
		{
			name: "valid error response",
			response: DashboardResponse{
				UserID:       1,
				TotalClicks:  0,
				TotalLinks:   0,
				RecentClicks: []RecentClick{},
				TopLinks:     []TopLink{},
				Status:       "error",
				Message:      &validMessage,
			},
			wantErr: false,
		},
		{
			name: "valid limited response",
			response: DashboardResponse{
				UserID:       1,
				TotalClicks:  0,
				TotalLinks:   0,
				RecentClicks: []RecentClick{},
				TopLinks:     []TopLink{},
				Status:       "limited",
			},
			wantErr: false,
		},
		{
			name: "invalid - zero user_id",
			response: DashboardResponse{
				UserID: 0,
				Status: "success",
			},
			wantErr: true,
		},
		{
			name: "invalid - negative total_clicks",
			response: DashboardResponse{
				UserID:      1,
				TotalClicks: -1,
				Status:      "success",
			},
			wantErr: true,
		},
		{
			name: "invalid - negative total_links",
			response: DashboardResponse{
				UserID:     1,
				TotalLinks: -1,
				Status:     "success",
			},
			wantErr: true,
		},
		{
			name: "invalid - unknown status",
			response: DashboardResponse{
				UserID: 1,
				Status: "unknown",
			},
			wantErr: true,
		},
		{
			name: "invalid - too many recent clicks",
			response: DashboardResponse{
				UserID: 1,
				RecentClicks: []RecentClick{
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, // 11 items
				},
				Status: "success",
			},
			wantErr: true,
		},
		{
			name: "invalid - too many top links",
			response: DashboardResponse{
				UserID: 1,
				TopLinks: []TopLink{
					{}, {}, {}, {}, {}, {}, // 6 items
				},
				Status: "success",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.response.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DashboardResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDashboardResponse_StatusCheckers(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		isSuccess  bool
		isError    bool
		isLimited  bool
	}{
		{
			name:      "success status",
			status:    "success",
			isSuccess: true,
			isError:   false,
			isLimited: false,
		},
		{
			name:      "error status",
			status:    "error",
			isSuccess: false,
			isError:   true,
			isLimited: false,
		},
		{
			name:      "limited status",
			status:    "limited",
			isSuccess: false,
			isError:   false,
			isLimited: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := DashboardResponse{
				UserID: 1,
				Status: tt.status,
			}

			if response.IsSuccess() != tt.isSuccess {
				t.Errorf("IsSuccess() = %v, want %v", response.IsSuccess(), tt.isSuccess)
			}
			if response.IsError() != tt.isError {
				t.Errorf("IsError() = %v, want %v", response.IsError(), tt.isError)
			}
			if response.IsLimited() != tt.isLimited {
				t.Errorf("IsLimited() = %v, want %v", response.IsLimited(), tt.isLimited)
			}
		})
	}
}

func TestDashboardResponse_GetMessage(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		msg := "test message"
		response := DashboardResponse{
			UserID:  1,
			Status:  "error",
			Message: &msg,
		}

		if response.GetMessage() != msg {
			t.Errorf("GetMessage() = %s, want %s", response.GetMessage(), msg)
		}
	})

	t.Run("without message", func(t *testing.T) {
		response := DashboardResponse{
			UserID: 1,
			Status: "success",
		}

		if response.GetMessage() != "" {
			t.Errorf("GetMessage() = %s, want empty string", response.GetMessage())
		}
	})
}

func TestDashboardResponse_JSON_Compatibility(t *testing.T) {
	// This simulates the JSON response from Python analytic-service
	pythonJSON := `{
		"user_id": 1,
		"total_clicks": 1542,
		"total_links": 45,
		"recent_clicks": [
			{
				"short_code": "my-link",
				"clicked_at": "2025-10-20T10:00:00Z",
				"ip_address_hash": "abc123hash",
				"user_agent": "Mozilla/5.0",
				"country_code": "ID",
				"city": "Jakarta",
				"device_type": "desktop",
				"browser_name": "Chrome",
				"is_bot": false
			}
		],
		"top_links": [
			{
				"short_code": "popular-link",
				"clicks": 350
			}
		],
		"status": "success"
	}`

	var response DashboardResponse
	if err := json.Unmarshal([]byte(pythonJSON), &response); err != nil {
		t.Fatalf("Failed to unmarshal Python JSON: %v", err)
	}

	// Validate the response
	if err := response.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Check values
	if response.UserID != 1 {
		t.Errorf("UserID = %d, want 1", response.UserID)
	}
	if response.TotalClicks != 1542 {
		t.Errorf("TotalClicks = %d, want 1542", response.TotalClicks)
	}
	if response.TotalLinks != 45 {
		t.Errorf("TotalLinks = %d, want 45", response.TotalLinks)
	}
	if !response.IsSuccess() {
		t.Errorf("IsSuccess() = false, want true")
	}
	if len(response.RecentClicks) != 1 {
		t.Errorf("len(RecentClicks) = %d, want 1", len(response.RecentClicks))
	}
	if len(response.TopLinks) != 1 {
		t.Errorf("len(TopLinks) = %d, want 1", len(response.TopLinks))
	}

	// Check nested data
	if response.RecentClicks[0].ShortCode != "my-link" {
		t.Errorf("RecentClicks[0].ShortCode = %s, want my-link", response.RecentClicks[0].ShortCode)
	}
	if response.TopLinks[0].Clicks != 350 {
		t.Errorf("TopLinks[0].Clicks = %d, want 350", response.TopLinks[0].Clicks)
	}
}

func TestDashboardResponse_JSON_ErrorResponse(t *testing.T) {
	// Simulate error response from Python
	pythonJSON := `{
		"user_id": 1,
		"status": "error",
		"message": "Database connection failed",
		"total_clicks": 0,
		"total_links": 0,
		"recent_clicks": [],
		"top_links": []
	}`

	var response DashboardResponse
	if err := json.Unmarshal([]byte(pythonJSON), &response); err != nil {
		t.Fatalf("Failed to unmarshal Python JSON: %v", err)
	}

	// Validate
	if err := response.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Check status
	if !response.IsError() {
		t.Errorf("IsError() = false, want true")
	}

	// Check message
	expectedMsg := "Database connection failed"
	if response.GetMessage() != expectedMsg {
		t.Errorf("GetMessage() = %s, want %s", response.GetMessage(), expectedMsg)
	}
}
