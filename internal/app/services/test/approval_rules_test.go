package services_test

import (
	"testing"

	"rule-based-approval-engine/internal/app/services/helpers"
)

func TestValidatePendingStatus(t *testing.T) {
	tests := []struct {
		status      string
		expectError bool
	}{
		{"PENDING", false},
		{"APPROVED", true},
		{"REJECTED", true},
	}

	for _, tt := range tests {
		err := helpers.ValidatePendingStatus(tt.status)

		if tt.expectError && err == nil {
			t.Errorf("expected error for status %s", tt.status)
		}

		if !tt.expectError && err != nil {
			t.Errorf("did not expect error for status %s", tt.status)
		}
	}
}
func TestValidateApproverRole(t *testing.T) {
	tests := []struct {
		name          string
		approverRole  string
		requesterRole string
		expectError   bool
	}{
		{"manager approves employee", "MANAGER", "EMPLOYEE", false},
		{"manager approves admin", "MANAGER", "ADMIN", true},
		{"admin approves manager", "ADMIN", "MANAGER", false},
		{"admin approves employee", "ADMIN", "EMPLOYEE", false},
		{"invalid approver", "EMPLOYEE", "EMPLOYEE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helpers.ValidateApproverRole(
				tt.approverRole,
				tt.requesterRole,
			)

			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
