package services

import (
	"context"
	"errors"
	"strings"

	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
)

func ApplyExpense(
	userID int64,
	amount float64,
	category string,
	reason string,
) (string, error) {

	ctx := context.Background()

	// ---- Input validations ----
	if userID <= 0 {
		return "", errors.New("invalid user")
	}

	if amount <= 0 {
		return "", apperrors.ErrInvalidExpenseAmount
	}

	if strings.TrimSpace(category) == "" {
		return "", apperrors.ErrInvalidExpenseCategory
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return "", errors.New("unable to start transaction")
	}
	defer tx.Rollback(ctx)

	// ---- Fetch remaining expense balance ----
	var remaining float64
	err = tx.QueryRow(
		ctx,
		`SELECT remaining_amount FROM expense WHERE user_id=$1`,
		userID,
	).Scan(&remaining)

	if err == pgx.ErrNoRows {
		return "", apperrors.ErrExpenseBalanceMissing
	}
	if err != nil {
		return "", errors.New("failed to fetch expense balance")
	}

	if amount > remaining {
		return "", apperrors.ErrExpenseLimitExceeded
	}

	// ---- Fetch user grade ----
	gradeID, err := FetchUserGrade(ctx, tx, userID)
	if err != nil {
		return "", err
	}

	// ---- Fetch rule ----
	rule, err := GetRule("EXPENSE", gradeID)
	if err != nil {
		return "", apperrors.ErrRuleNotFound
	}

	// ---- Decision ----
	result := MakeDecision("EXPENSE", rule.Condition, amount)
	status := result.Status
	message := result.Message

	// ---- Insert expense request ----
	_, err = tx.Exec(
		ctx,
		`INSERT INTO expense_requests
		 (employee_id, amount, category, reason, status, rule_id)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		userID, amount, category, reason, status, rule.ID,
	)
	if err != nil {
		return "", errors.New("failed to create expense request")
	}

	// ---- Deduct balance if auto-approved ----
	if status == "AUTO_APPROVED" {
		_, err = tx.Exec(
			ctx,
			`UPDATE expense
			 SET remaining_amount = remaining_amount - $1
			 WHERE user_id=$2`,
			amount, userID,
		)
		if err != nil {
			return "", errors.New("failed to update expense balance")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", errors.New("failed to commit transaction")
	}

	return message, nil
}

func CancelExpense(userID, requestID int64) error {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var status string
	var amount float64

	err = tx.QueryRow(
		ctx,
		`SELECT status, amount
		 FROM expense_requests
		 WHERE id=$1 AND employee_id=$2`,
		requestID, userID,
	).Scan(&status, &amount)

	if err != nil {
		return err
	}
	// reuse CanCancel from apply_cancel_rules.go
	if err := CanCancel(status); err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE expense_requests
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
			`UPDATE expense
			 SET remaining_amount = remaining_amount + $1
			 WHERE user_id=$2`,
			amount, userID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
