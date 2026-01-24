package services

import (
	"context"
	"errors"

	"rule-based-approval-engine/internal/apperrors"
	"rule-based-approval-engine/internal/database"

	"github.com/jackc/pgx/v5"
)

func ApplyDiscount(
	userID int64,
	percent float64,
	reason string,
) (string, error) {

	ctx := context.Background()

	// ---- Input validations ----
	if userID <= 0 {
		return "", errors.New("invalid user")
	}

	if percent <= 0 {
		return "", apperrors.ErrInvalidDiscountPercent
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return "", errors.New("unable to start transaction")
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
		return "", apperrors.ErrDiscountBalanceMissing
	}
	if err != nil {
		return "", errors.New("failed to fetch discount balance")
	}

	if percent > remaining {
		return "", apperrors.ErrDiscountLimitExceeded
	}

	// ---- Fetch user grade ----
	var gradeID int64
	err = tx.QueryRow(
		ctx,
		`SELECT grade_id FROM users WHERE id=$1`,
		userID,
	).Scan(&gradeID)

	if err == pgx.ErrNoRows {
		return "", apperrors.ErrUserNotFound
	}
	if err != nil {
		return "", errors.New("failed to fetch user grade")
	}

	// ---- Fetch rule ----
	rule, err := GetRule("DISCOUNT", gradeID)
	if err != nil {
		return "", apperrors.ErrRuleNotFound
	}

	// ---- Decision ----
	decision := Decide("DISCOUNT", rule.Condition, percent)

	status := "PENDING"
	message := "Discount submitted to manager for approval"

	if decision == "AUTO_APPROVE" {
		status = "AUTO_APPROVED"
		message = "Discount approved by system"
	}

	// ---- Insert discount request ----
	_, err = tx.Exec(
		ctx,
		`INSERT INTO discount_requests
		 (employee_id, discount_percentage, reason, status, rule_id)
		 VALUES ($1,$2,$3,$4,$5)`,
		userID, percent, reason, status, rule.ID,
	)
	if err != nil {
		return "", errors.New("failed to create discount request")
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
			return "", errors.New("failed to update discount balance")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", errors.New("failed to commit transaction")
	}

	return message, nil
}
func CancelDiscount(userID, requestID int64) error {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return errors.New("unable to start transaction")
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
		return errors.New("failed to fetch discount request")
	}

	// Cannot cancel finalized request
	if status == "APPROVED" || status == "REJECTED" || status == "CANCELLED" {
		return apperrors.ErrDiscountCannotCancel
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE discount_requests
		 SET status='CANCELLED'
		 WHERE id=$1`,
		requestID,
	)
	if err != nil {
		return errors.New("failed to cancel discount request")
	}

	// 🔄 Restore discount if auto-approved
	if status == "AUTO_APPROVED" {
		_, err = tx.Exec(
			ctx,
			`UPDATE discount
			 SET remaining_discount = remaining_discount + $1
			 WHERE user_id=$2`,
			percent, userID,
		)
		if err != nil {
			return errors.New("failed to restore discount balance")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}
