package services

import (
	"context"
	"errors"
	"time"

	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/pkg/helpers"

	"github.com/jackc/pgx/v5"
)

func GetPendingExpenseRequests(role string, approverID int64) ([]map[string]interface{}, error) {
	ctx := context.Background()

	var rows pgx.Rows
	var err error

	if role == "MANAGER" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT er.id, u.name, er.amount, er.category , er.reason,er.created_at 
			 FROM expense_requests er
			 JOIN users u ON er.employee_id = u.id
			 WHERE er.status='PENDING' AND u.manager_id=$1`,
			approverID,
		)
	} else if role == "ADMIN" {
		rows, err = database.DB.Query(
			ctx,
			`SELECT er.id, u.name, er.amount, er.category , er.reason,er.created_at
			 FROM expense_requests er
			 JOIN users u ON er.employee_id = u.id
			 WHERE er.status='PENDING'`,
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
		var reason *string
		var amount float64
		var createdAt time.Time
		if err := rows.Scan(&id, &name, &amount, &category, &reason, &createdAt); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":         id,
			"employee":   name,
			"amount":     amount,
			"category":   category,
			"reason":     reason,
			"created_at": createdAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
func ApproveExpense(
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

	// Fetch request
	err = tx.QueryRow(
		ctx,
		`SELECT employee_id, status
		 FROM expense_requests
		 WHERE id=$1`,
		requestID,
	).Scan(&employeeID, &status)

	if err != nil {
		return err
	}

	//  Validate pending
	if err := helpers.ValidatePendingStatus(status); err != nil {
		return err
	}
	if approverID == employeeID {
		return errors.New("self approval is not allowed")
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

	// 4️⃣ Validate authorization
	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// 5️⃣ Default comment
	if comment == "" {
		comment = "Approved"
	}

	// 6️⃣ Update request
	_, err = tx.Exec(
		ctx,
		`UPDATE expense_requests
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

func RejectExpense(
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

	// 1. Fetch request
	err = tx.QueryRow(
		ctx,
		`SELECT employee_id, status
		 FROM expense_requests
		 WHERE id=$1`,
		requestID,
	).Scan(&employeeID, &status)

	if err != nil {
		return err
	}

	// 2. Validate pending status
	if err := helpers.ValidatePendingStatus(status); err != nil {
		return err
	}

	// 3. Prevent self-approval
	if approverID == employeeID {
		return errors.New("self approval is not allowed")
	}

	// 4. Fetch requester role
	var requesterRole string
	err = tx.QueryRow(
		ctx,
		`SELECT role FROM users WHERE id=$1`,
		employeeID,
	).Scan(&requesterRole)

	if err != nil {
		return err
	}

	// 5. Validate approver role
	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// 6. Default rejection comment
	if comment == "" {
		comment = "Rejected"
	}

	// 7. Update request
	_, err = tx.Exec(
		ctx,
		`UPDATE expense_requests
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
