package services

import (
	"context"
	"errors"

	"rule-based-approval-engine/internal/database"

	"github.com/jackc/pgx/v5"
)

func GetPendingExpenseRequests(role string, approverID int64) ([]map[string]interface{}, error) {
	ctx := context.Background()

	var rows pgx.Rows
	var err error

	if role == "MANAGER" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT er.id, u.name, er.amount, er.category
			 FROM expense_requests er
			 JOIN users u ON er.employee_id = u.id
			 WHERE er.status='PENDING' AND u.manager_id=$1`,
			approverID,
		)
	} else if role == "ADMIN" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT er.id, u.name, er.amount, er.category
			 FROM expense_requests er
			 JOIN users u ON er.employee_id = u.id
			 WHERE er.status='PENDING' AND u.role='MANAGER'`,
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
		var name, category string
		var amount float64

		if err := rows.Scan(&id, &name, &amount, &category); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":       id,
			"employee": name,
			"amount":   amount,
			"category": category,
		})
	}

	return result, nil
}

func ApproveExpense(role string, approverID, requestID int64) error {
	ctx := context.Background()

	_, err := database.DB.Exec(
		ctx,
		`UPDATE expense_requests SET status='APPROVED'
		 WHERE id=$1 AND status='PENDING'`,
		requestID,
	)
	return err
}

func RejectExpense(role string, approverID, requestID int64) error {
	ctx := context.Background()

	_, err := database.DB.Exec(
		ctx,
		`UPDATE expense_requests SET status='REJECTED'
		 WHERE id=$1 AND status='PENDING'`,
		requestID,
	)
	return err
}
