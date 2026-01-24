package services

import (
	"context"

	"rule-based-approval-engine/internal/apperrors"
	"rule-based-approval-engine/internal/database"

	"github.com/jackc/pgx/v5"
)

func GetPendingDiscountRequests(role string, approverID int64) ([]map[string]interface{}, error) {
	ctx := context.Background()

	var rows pgx.Rows
	var err error

	if role == "MANAGER" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT dr.id, u.name, dr.discount_percentage
			 FROM discount_requests dr
			 JOIN users u ON dr.employee_id = u.id
			 WHERE dr.status='PENDING' AND u.manager_id=$1`,
			approverID,
		)
	} else if role == "ADMIN" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT dr.id, u.name, dr.discount_percentage
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

		if err := rows.Scan(&id, &name, &percentage); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":         id,
			"employee":   name,
			"percentage": percentage,
		})
	}

	return result, nil
}

func ApproveDiscount(role string, approverID, requestID int64, comment string) error {
	return processDiscountApproval(role, approverID, requestID, comment, "APPROVED")
}

func RejectDiscount(role string, approverID, requestID int64, comment string) error {
	return processDiscountApproval(role, approverID, requestID, comment, "REJECTED")
}

func processDiscountApproval(
	role string,
	approverID, requestID int64,
	comment string,
	newStatus string,
) error {

	ctx := context.Background()

	if role != "MANAGER" && role != "ADMIN" {
		return apperrors.ErrUnauthorizedApprover
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var employeeID int64
	var status string

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

	if status != "PENDING" {
		return apperrors.ErrDiscountRequestNotPending
	}

	var requesterRole string
	err = tx.QueryRow(
		ctx,
		`SELECT role FROM users WHERE id=$1`,
		employeeID,
	).Scan(&requesterRole)

	if err != nil {
		return err
	}

	switch role {
	case "MANAGER":
		if requesterRole != "EMPLOYEE" {
			return apperrors.ErrUnauthorizedApprover
		}
	case "ADMIN":
		if requesterRole != "EMPLOYEE" && requesterRole != "MANAGER" {
			return apperrors.ErrUnauthorizedApprover
		}
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE discount_requests
		 SET status=$1,
		     approved_by_id=$2,
		     approval_comment=$3
		 WHERE id=$4`,
		newStatus, approverID, comment, requestID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
