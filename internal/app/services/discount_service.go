package services

import (
	"context"

	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
)

func ApplyDiscount(
	userID int64,
	percent float64,
	reason string,
) (string, string, error) {

	ctx := context.Background()

	// ---- Input validations ----
	if userID <= 0 {
		return "", "", apperrors.ErrInvalidUser
	}

	if percent <= 0 {
		return "", "", apperrors.ErrInvalidDiscountPercent
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return "", "", apperrors.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	// ---- Fetch remaining discount ----
	var remaining float64
	err = tx.QueryRow(
		ctx,
		`SELECT remaining_discount FROM discount WHERE user_id=$1`,
		userID,
	).Scan(&remaining)

	if err == pgx.ErrNoRows {
		return "", "", apperrors.ErrDiscountBalanceMissing
	}
	if err != nil {
		return "", "", apperrors.ErrBalanceFetchFailed
	}

	if percent > remaining {
		return "", "", apperrors.ErrDiscountLimitExceeded
	}

	// ---- Fetch user grade ----
	gradeID, err := helpers.FetchUserGrade(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	// ---- Fetch rule ----
	rule, err := GetRule("DISCOUNT", gradeID)
	if err != nil {
		return "", "", apperrors.ErrRuleNotFound
	}

	// ---- Decision ----
	result := helpers.MakeDecision("DISCOUNT", rule.Condition, percent)
	status := result.Status
	message := result.Message

	// ---- Insert discount request ----
	_, err = tx.Exec(
		ctx,
		`INSERT INTO discount_requests
		 (employee_id, discount_percentage, reason, status, rule_id)
		 VALUES ($1,$2,$3,$4,$5)`,
		userID, percent, reason, status, rule.ID,
	)
	if err != nil {
		return "", "", apperrors.ErrInsertFailed
	}

	// ---- Deduct discount if auto-approved ----
	if status == "AUTO_APPROVED" {
		_, err = tx.Exec(
			ctx,
			`UPDATE discount
			 SET remaining_discount = remaining_discount - $1
			 WHERE user_id=$2`,
			percent, userID,
		)
		if err != nil {
			return "", "", apperrors.ErrBalanceUpdateFailed
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", apperrors.ErrTransactionCommit
	}

	return message, status, nil
}
func CancelDiscount(userID, requestID int64) error {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return apperrors.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	var status string
	var percent float64

	err = tx.QueryRow(
		ctx,
		`SELECT status, discount_percentage
		 FROM discount_requests
		 WHERE id=$1 AND employee_id=$2`,
		requestID, userID,
	).Scan(&status, &percent)

	//  HANDLE NO ROWS PROPERLY
	if err == pgx.ErrNoRows {
		return apperrors.ErrDiscountRequestNotFound
	}
	if err != nil {
		return apperrors.ErrQueryFailed
	}

	// reuse CanCancel from apply_cancel_rules.go
	if err := helpers.CanCancel(status); err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE discount_requests
		 SET status='CANCELLED'
		 WHERE id=$1`,
		requestID,
	)
	if err != nil {
		return apperrors.ErrUpdateFailed
	}

	//  Restore discount if auto-approved
	if status == "AUTO_APPROVED" {
		_, err = tx.Exec(
			ctx,
			`UPDATE discount
			 SET remaining_discount = remaining_discount + $1
			 WHERE user_id=$2`,
			percent, userID,
		)
		if err != nil {
			return apperrors.ErrBalanceUpdateFailed
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.ErrTransactionCommit
	}

	return nil
}
