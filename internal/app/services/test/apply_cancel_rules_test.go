package services_test

import (
	"rule-based-approval-engine/internal/app/services/helpers"
	"testing"
)

func TestMakeDecision(t *testing.T) {
	tests := []struct {
		name         string
		requestType  string
		condition    map[string]interface{}
		value        float64
		expectStatus string
	}{
		{
			name:        "auto approve leave",
			requestType: "LEAVE",
			condition: map[string]interface{}{
				"max_days": float64(3),
			},
			value:        2,
			expectStatus: "AUTO_APPROVED",
		},
		{
			name:        "pending leave",
			requestType: "LEAVE",
			condition: map[string]interface{}{
				"max_days": float64(3),
			},
			value:        15,
			expectStatus: "PENDING",
		},
		{
			name:        "auto approve expense",
			requestType: "EXPENSE",
			condition: map[string]interface{}{
				"max_amount": float64(5000),
			},
			value:        3000,
			expectStatus: "AUTO_APPROVED",
		},
		{
			name:        "pending expense",
			requestType: "EXPENSE",
			condition: map[string]interface{}{
				"max_amount": float64(5000),
			},
			value:        8000,
			expectStatus: "PENDING",
		},
		{
			name:        "auto approve discount",
			requestType: "DISCOUNT",
			condition: map[string]interface{}{
				"max_percent": float64(10),
			},
			value:        8,
			expectStatus: "AUTO_APPROVED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.MakeDecision(tt.requestType, tt.condition, tt.value)

			if result.Status != tt.expectStatus {
				t.Errorf(
					"expected status %s, got %s",
					tt.expectStatus,
					result.Status,
				)
			}

			if result.Message == "" {
				t.Errorf("expected non-empty message")
			}
		})
	}
}

func TestCanCancel(t *testing.T) {
	tests := []struct {
		status      string
		expectError bool
	}{
		{"PENDING", false},
		{"AUTO_APPROVED", false},
		{"APPROVED", true},
		{"REJECTED", true},
		{"CANCELLED", true},
	}

	for _, tt := range tests {
		err := helpers.CanCancel(tt.status)

		if tt.expectError && err == nil {
			t.Errorf("expected error for status %s, got nil", tt.status)
		}

		if !tt.expectError && err != nil {
			t.Errorf("did not expect error for status %s, got %v", tt.status, err)
		}
	}
}
