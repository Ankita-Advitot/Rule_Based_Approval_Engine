package services

import (
	"context"
	"time"

	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/utils"

	"github.com/jackc/pgx/v5"
)

func ApplyLeave(
	userID int64,
	from time.Time,
	to time.Time,
	days int,
	leaveType string,
	reason string,
) (string, string, error) {
	ctx := context.Background()

	// ---- Input validations ----
	if userID <= 0 {
		return "", "", apperrors.ErrInvalidUser
	}

	if days <= 0 {
		return "", "", apperrors.ErrInvalidLeaveDays
	}

	if from.After(to) {
		return "", "", apperrors.ErrInvalidDateRange
	}
	// ---- Overlap validation ----
	overlap, err := HasOverlappingLeave(ctx, userID, from, to)
	if err != nil {
		return "", "", apperrors.ErrLeaveVerificationFailed
	}

	if overlap {
		return "", "", apperrors.ErrLeaveOverlap
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return "", "", apperrors.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	// ---- Fetch remaining leave balance ----
	var remaining int
	err = tx.QueryRow(
		ctx,
		`SELECT remaining_count FROM leaves WHERE user_id=$1`,
		userID,
	).Scan(&remaining)

	if err == pgx.ErrNoRows {
		return "", "", apperrors.ErrLeaveBalanceMissing
	}
	if err != nil {
		return "", "", apperrors.ErrBalanceFetchFailed
	}

	if days > remaining {
		return "", "", apperrors.ErrLeaveBalanceExceeded
	}

	// ---- Fetch user grade ----
	gradeID, err := helpers.FetchUserGrade(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	// ---- Fetch rule ----
	rule, err := GetRule("LEAVE", gradeID)
	if err != nil {
		return "", "", apperrors.ErrRuleNotFound
	}

	// ---- Decision ----
	result := helpers.MakeDecision("LEAVE", rule.Condition, float64(days))
	status := result.Status
	message := result.Message

	_, err = tx.Exec(
		ctx,
		`INSERT INTO leave_requests
		 (employee_id, from_date, to_date, reason, leave_type, status, rule_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		userID, from, to, reason, leaveType, status, rule.ID,
	)

	if err != nil {
		return "", "", helpers.MapPgError(err)
	}

	// ---- Deduct balance if auto-approved ----
	if status == "AUTO_APPROVED" {
		_, err = tx.Exec(
			ctx,
			`UPDATE leaves
			 SET remaining_count = remaining_count - $1
			 WHERE user_id=$2`,
			days, userID,
		)
		if err != nil {
			return "", "", helpers.MapPgError(err)
		}

	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", apperrors.ErrTransactionCommit
	}

	return message, status, nil
}

func HasOverlappingLeave(
	ctx context.Context,
	userID int64,
	fromDate, toDate time.Time,
) (bool, error) {

	var dummy int

	err := database.DB.QueryRow(
		ctx,
		`SELECT 1
		 FROM leave_requests
		 WHERE employee_id = $1
		   AND status IN ('PENDING', 'APPROVED', 'AUTO_APPROVED') 
		   AND from_date <= $2
		   AND to_date >= $3
		 LIMIT 1`,
		userID,
		toDate,
		fromDate,
	).Scan(&dummy)

	// pgx NO ROWS = no overlap
	if err == pgx.ErrNoRows {
		return false, nil
	}

	//  real system error
	if err != nil {
		return false, err
	}

	//  overlap exists
	return true, nil
}

func CancelLeave(userID, requestID int64) error {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var status string
	var from, to time.Time

	err = tx.QueryRow(
		ctx,
		`SELECT status, from_date, to_date 
		 FROM leave_requests 
		 WHERE id=$1 AND employee_id=$2`,
		requestID, userID,
	).Scan(&status, &from, &to)

	if err != nil {
		return err
	}

	// reuse CanCancel from apply_cancel_rules.go
	if err := helpers.CanCancel(status); err != nil {
		return err
	}

	days := utils.CalculateLeaveDays(from, to)

	_, err = tx.Exec(
		ctx,
		`UPDATE leave_requests 
		 SET status='CANCELLED' 
		 WHERE id=$1`,
		requestID,
	)
	if err != nil {
		return err
	}

	if status == "AUTO_APPROVED" {
		_, err = tx.Exec(
			ctx,
			`UPDATE leaves 
			 SET remaining_count = remaining_count + $1
			 WHERE user_id=$2`,
			days, userID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
