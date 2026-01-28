package services

import (
	"context"
	"errors"
	"fmt"
	"rule-based-approval-engine/internal/apperrors"

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
