package helpers

import (
	"rule-based-approval-engine/internal/pkg/apperrors"
)

func ValidatePendingStatus(status string) error {
	if status != "PENDING" {
		return apperrors.ErrRequestNotPending
	}
	return nil
}

func ValidateApproverRole(approverRole, requesterRole string) error {

	// Employee can never approve
	if approverRole == "EMPLOYEE" {
		return apperrors.ErrEmployeeCannotApprove
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
			return apperrors.ErrManagerNeedsAdmin
		}

		// Admin never requests, but safety net
		if requesterRole == "ADMIN" {
			return apperrors.ErrAdminRequestNotAllowed
		}
	}

	return apperrors.ErrUnauthorizedApproval
}
