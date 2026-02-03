package helpers

import (
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/pkg/apperrors"
)

func ValidatePendingStatus(status string) error {
	if status != constants.StatusPending {
		return apperrors.ErrRequestNotPending
	}
	return nil
}

func ValidateApproverRole(approverRole, requesterRole string) error {

	// Employee can never approve
	if approverRole == constants.RoleEmployee {
		return apperrors.ErrEmployeeCannotApprove
	}

	// Admin can approve anyone
	if approverRole == constants.RoleAdmin {
		return nil
	}

	// Manager rules
	if approverRole == constants.RoleManager {

		// Manager can approve employee only
		if requesterRole == constants.RoleEmployee {
			return nil
		}

		// Manager approving manager
		if requesterRole == constants.RoleManager {
			return apperrors.ErrManagerNeedsAdmin
		}

		// Admin never requests, but safety net
		if requesterRole == constants.RoleAdmin {
			return apperrors.ErrAdminRequestNotAllowed
		}
	}

	return apperrors.ErrUnauthorizedApproval
}
