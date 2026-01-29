package services

import (
	"context"
	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/pkg/utils"
	"time"
)

func AutoRejectLeaveRequests() error {
	ctx := context.Background()

	rows, err := database.DB.Query(
		ctx,
		`SELECT id, created_at 
		 FROM leave_requests 
		 WHERE status='PENDING'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	now := time.Now()

	for rows.Next() {
		var id int64
		var createdAt time.Time

		if err := rows.Scan(&id, &createdAt); err != nil {
			return err
		}

		workingDays := utils.CountWorkingDays(createdAt, now)

		if workingDays >= 7 {
			_, err = database.DB.Exec(
				ctx,
				`UPDATE leave_requests
				 SET status='AUTO_REJECTED',
				     approval_comment='Auto rejected after 7 working days'
				 WHERE id=$1`,
				id,
			)
			if err != nil {
				return err
			}
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
func AutoRejectExpenseRequests() error {
	ctx := context.Background()

	rows, err := database.DB.Query(
		ctx,
		`SELECT id, created_at 
		 FROM expense_requests 
		 WHERE status='PENDING'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	now := time.Now()

	for rows.Next() {
		var id int64
		var createdAt time.Time

		if err := rows.Scan(&id, &createdAt); err != nil {
			return err
		}

		workingDays := utils.CountWorkingDays(createdAt, now)

		if workingDays >= 7 {
			_, err = database.DB.Exec(
				ctx,
				`UPDATE expense_requests
				 SET status='AUTO_REJECTED',
				     approval_comment='Auto rejected after 7 working days'
				 WHERE id=$1`,
				id,
			)
			if err != nil {
				return err
			}
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
func AutoRejectDiscountRequests() error {
	ctx := context.Background()

	rows, err := database.DB.Query(
		ctx,
		`SELECT id, created_at 
		 FROM discount_requests 
		 WHERE status='PENDING'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	now := time.Now()

	for rows.Next() {
		var id int64
		var createdAt time.Time

		if err := rows.Scan(&id, &createdAt); err != nil {
			return err
		}

		workingDays := utils.CountWorkingDays(createdAt, now)

		if workingDays >= 7 {
			_, err = database.DB.Exec(
				ctx,
				`UPDATE discount_requests
				 SET status='AUTO_REJECTED',
				     approval_comment='Auto rejected after 7 working days'
				 WHERE id=$1`,
				id,
			)
			if err != nil {
				return err
			}
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
