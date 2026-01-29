package services

import (
	"context"
	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
)

func InitializeBalances(tx pgx.Tx, userID int64, gradeID int64) error {
	ctx := context.Background()

	var leaveLimit int
	var expenseLimit float64
	var discountLimit float64

	err := tx.QueryRow(
		ctx,
		`SELECT annual_leave_limit, annual_expense_limit, discount_limit_percent
		 FROM grades WHERE id=$1`,
		gradeID,
	).Scan(&leaveLimit, &expenseLimit, &discountLimit)

	if err != nil {
		if err == pgx.ErrNoRows {
			return apperrors.ErrQueryFailed // Or a more specific error if grade not found
		}
		return helpers.MapPgError(err)
	}

	// Leave wallet
	_, err = tx.Exec(
		ctx,
		`INSERT INTO leaves (user_id, total_allocated, remaining_count)
		 VALUES ($1,$2,$2)
		 ON CONFLICT (user_id) DO NOTHING`,
		userID, leaveLimit,
	)
	if err != nil {
		return helpers.MapPgError(err)
	}

	// Expense wallet
	_, err = tx.Exec(
		ctx,
		`INSERT INTO expense (user_id, total_amount, remaining_amount)
		 VALUES ($1,$2,$2)
		 ON CONFLICT (user_id) DO NOTHING`,
		userID, expenseLimit,
	)
	if err != nil {
		return helpers.MapPgError(err)
	}

	// Discount wallet
	_, err = tx.Exec(
		ctx,
		`INSERT INTO discount (user_id, total_discount, remaining_discount)
		 VALUES ($1,$2,$2)
		 ON CONFLICT (user_id) DO NOTHING`,
		userID, discountLimit,
	)
	if err != nil {
		return helpers.MapPgError(err)
	}

	return nil
}
