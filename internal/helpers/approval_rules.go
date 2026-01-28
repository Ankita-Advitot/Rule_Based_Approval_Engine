package helpers

import (
	"errors"
)

func ValidatePendingStatus(status string) error {
	if status != "PENDING" {
		return errors.New("request not pending")
	}
	return nil
}

func ValidateApproverRole(approverRole, requesterRole string) error {

	// Employee can never approve
	if approverRole == "EMPLOYEE" {
		return errors.New("employees are not allowed to approve requests")
	}

	// Admin can approve anyone
	if approverRole == "ADMIN" {
		return nil
	}

	// Manager rules
	if approverRole == "MANAGER" {

		// Manager can approve employee only
		if requesterRole == "EMPLOYEE" {
			return nil
		}

		// Manager approving manager
		if requesterRole == "MANAGER" {
			return errors.New("managers can only be approved by admin")
		}

		// Admin never requests, but safety net
		if requesterRole == "ADMIN" {
			return errors.New("admin requests are not allowed")
		}
	}

	return errors.New("unauthorized approval attempt")
}
