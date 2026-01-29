package helpers

import (
	"context"
	"errors"
	"fmt"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
)

type DecisionResult struct {
	Status  string
	Message string
}

func MakeDecision(
	requestType string,
	condition map[string]interface{},
	value float64,
) DecisionResult {

	decision := Decide(requestType, condition, value)
	fmt.Println("DEBUG Decide returned:", decision)

	if decision == "AUTO_APPROVE" {
		return DecisionResult{
			Status:  "AUTO_APPROVED",
			Message: requestType + " approved by system",
		}
	}

	return DecisionResult{
		Status:  "PENDING",
		Message: requestType + " submitted for approval",
	}
}

func CanCancel(status string) error {
	switch status {
	case "APPROVED", "REJECTED", "CANCELLED":
		return errors.New("cannot cancel finalized request")
	default:
		return nil
	}
}
func FetchUserGrade(ctx context.Context, tx pgx.Tx, userID int64) (int64, error) {
	var gradeID int64
	err := tx.QueryRow(
		ctx,
		`SELECT grade_id FROM users WHERE id=$1`,
		userID,
	).Scan(&gradeID)

	if err == pgx.ErrNoRows {
		return 0, apperrors.ErrUserNotFound
	}
	return gradeID, err
}
func Decide(
	requestType string,
	rule map[string]interface{},
	value float64,
) string {

	switch requestType {
	case "LEAVE":
		if EvaluateLeaveRule(rule, int(value)) {
			return "AUTO_APPROVE"
		}
	case "EXPENSE":
		if EvaluateExpenseRule(rule, value) {
			return "AUTO_APPROVE"
		}
	case "DISCOUNT":
		if EvaluateDiscountRule(rule, value) {
			return "AUTO_APPROVE"
		}
	}

	return "MANUAL"
}

func EvaluateLeaveRule(rule map[string]interface{}, days int) bool {
	maxDays, ok := rule["max_days"].(float64)
	fmt.Println("EvaluateLeaveRule ", ok)
	if !ok {
		return false
	}
	return days <= int(maxDays)
}

func EvaluateExpenseRule(rule map[string]interface{}, amount float64) bool {
	maxAmount, ok := rule["max_amount"].(float64)
	if !ok {
		return false
	}
	return amount <= maxAmount
}

func EvaluateDiscountRule(rule map[string]interface{}, percent float64) bool {
	maxPercent, ok := rule["max_percent"].(float64)
	if !ok {
		return false
	}
	return percent <= maxPercent
}
