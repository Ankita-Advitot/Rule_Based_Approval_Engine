package services

import (
	"context"
	"errors"

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
			`SELECT dr.id, u.name, dr.percentage
			 FROM discount_requests dr
			 JOIN users u ON dr.employee_id = u.id
			 WHERE dr.status='PENDING' AND u.manager_id=$1`,
			approverID,
		)
	} else if role == "ADMIN" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT dr.id, u.name, dr.percentage
			 FROM discount_requests dr
			 JOIN users u ON dr.employee_id = u.id
			 WHERE dr.status='PENDING' AND u.role='MANAGER'`,
		)
	} else {
		return nil, errors.New("unauthorized")
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

func ApproveDiscount(role string, approverID, requestID int64) error {
	ctx := context.Background()

	_, err := database.DB.Exec(
		ctx,
		`UPDATE discount_requests SET status='APPROVED'
		 WHERE id=$1 AND status='PENDING'`,
		requestID,
	)
	return err
}

func RejectDiscount(role string, approverID, requestID int64) error {
	ctx := context.Background()

	_, err := database.DB.Exec(
		ctx,
		`UPDATE discount_requests SET status='REJECTED'
		 WHERE id=$1 AND status='PENDING'`,
		requestID,
	)
	return err
}
