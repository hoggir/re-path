package domain

import (
	"encoding/json"
	"fmt"
)

// ContractValidator provides utilities for validating RPC contracts
type ContractValidator struct{}

// ValidateDashboardRequestJSON validates a dashboard request from JSON
func (v *ContractValidator) ValidateDashboardRequestJSON(data []byte) (*DashboardRequest, error) {
	var request DashboardRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &request, nil
}

// ValidateDashboardResponseJSON validates a dashboard response from JSON
func (v *ContractValidator) ValidateDashboardResponseJSON(data []byte) (*DashboardResponse, error) {
	var response DashboardResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	if err := response.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &response, nil
}

// MarshalDashboardRequest marshals a dashboard request to JSON with validation
func (v *ContractValidator) MarshalDashboardRequest(request *DashboardRequest) ([]byte, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return json.Marshal(request)
}

// MarshalDashboardResponse marshals a dashboard response to JSON with validation
func (v *ContractValidator) MarshalDashboardResponse(response *DashboardResponse) ([]byte, error) {
	if err := response.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return json.Marshal(response)
}

// Global validator instance
var Validator = &ContractValidator{}
