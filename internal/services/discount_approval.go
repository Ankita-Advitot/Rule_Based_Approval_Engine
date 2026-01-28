package services

import (
	"context"
	"errors"
	"time"

	"rule-based-approval-engine/internal/apperrors"
	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/helpers"

	"github.com/jackc/pgx/v5"
)

func GetPendingDiscountRequests(role string, approverID int64) ([]map[string]interface{}, error) {
	ctx := context.Background()

	var rows pgx.Rows
	var err error

	if role == "MANAGER" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT dr.id, u.name, dr.discount_percentage ,dr.reason,dr.created_at 
			 FROM discount_requests dr
			 JOIN users u ON dr.employee_id = u.id
			 WHERE dr.status='PENDING' AND u.manager_id=$1`,
			approverID,
		)
	} else if role == "ADMIN" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT dr.id, u.name, dr.discount_percentage ,dr.reason,dr.created_at
			 FROM discount_requests dr
			 JOIN users u ON dr.employee_id = u.id
			 WHERE dr.status='PENDING'`,
		)
	} else {
		return nil, apperrors.ErrUnauthorizedApprover
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int64
		var name string
		var percentage float64
		var reason string
		var createdAt time.Time

		if err := rows.Scan(&id, &name, &percentage, &reason, &createdAt); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":         id,
			"employee":   name,
			"percentage": percentage,
			"reason":     reason,
			"created_at": createdAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
func ApproveDiscount(
	role string,
	approverID, requestID int64,
	comment string,
) error {

	ctx := context.Background()

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var employeeID int64
	var status string

	//  Fetch request
	err = tx.QueryRow(
		ctx,
		`SELECT employee_id, status
		 FROM discount_requests
		 WHERE id=$1`,
		requestID,
	).Scan(&employeeID, &status)

	if err == pgx.ErrNoRows {
		return apperrors.ErrDiscountRequestNotFound
	}
	if err != nil {
		return err
	}

	if approverID == employeeID {
		return errors.New("self approval is not allowed")
	}
	//  Validate pending
	if err := helpers.ValidatePendingStatus(status); err != nil {
		return err
	}

	//  Fetch requester role
	var requesterRole string
	err = tx.QueryRow(
		ctx,
		`SELECT role FROM users WHERE id=$1`,
		employeeID,
	).Scan(&requesterRole)

	if err != nil {
		return err
	}

	// Validate approver role
	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// Default comment
	if comment == "" {
		comment = "Approved"
	}

	// Update request
	_, err = tx.Exec(
		ctx,
		`UPDATE discount_requests
		 SET status='APPROVED',
		     approved_by_id=$1,
		     approval_comment=$2
		 WHERE id=$3`,
		approverID, comment, requestID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
func RejectDiscount(
	role string,
	approverID, requestID int64,
	comment string,
) error {

	ctx := context.Background()

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var employeeID int64
	var status string

	//  Fetch request
	err = tx.QueryRow(
		ctx,
		`SELECT employee_id, status
		 FROM discount_requests
		 WHERE id=$1`,
		requestID,
	).Scan(&employeeID, &status)

	if err == pgx.ErrNoRows {
		return apperrors.ErrDiscountRequestNotFound
	}
	if err != nil {
		return err
	}
	if approverID == employeeID {
		return errors.New("self approval is not allowed")
	}
	//  Validate pending
	if err := helpers.ValidatePendingStatus(status); err != nil {
		return err
	}

	//  Fetch requester role
	var requesterRole string
	err = tx.QueryRow(
		ctx,
		`SELECT role FROM users WHERE id=$1`,
		employeeID,
	).Scan(&requesterRole)

	if err != nil {
		return err
	}

	//  Validate approver role
	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	//  Default comment
	if comment == "" {
		comment = "Rejected"
	}

	//  Update request
	_, err = tx.Exec(
		ctx,
		`UPDATE discount_requests
		 SET status='REJECTED',
		     approved_by_id=$1,
		     approval_comment=$2
		 WHERE id=$3`,
		approverID, comment, requestID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
