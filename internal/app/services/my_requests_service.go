package services

import (
	"context"
	"fmt"
	"time"

	"rule-based-approval-engine/internal/database"
)

func GetMyLeaveRequests(userID int64) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, leave_type, from_date, to_date, status, reason, approval_comment, created_at
		 FROM leave_requests
		 WHERE employee_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var (
			id        int64
			leaveType string
			fromDate  time.Time
			toDate    time.Time
			status    string
			reason    string
			comment   *string
			createdAt time.Time
		)

		if err := rows.Scan(
			&id,
			&leaveType,
			&fromDate,
			&toDate,
			&status,
			&reason,
			&comment,
			&createdAt,
		); err != nil {
			return nil, err
		}

		// Validation / formatting
		response := map[string]interface{}{
			"id":         id,
			"leave_type": leaveType,
			"from_date":  fromDate.Format("2006-01-02"),
			"to_date":    toDate.Format("2006-01-02"),
			"status":     status,
			"reason":     reason,
			"created_at": createdAt.Format(time.RFC3339),
		}

		if comment != nil {
			response["approval_comment"] = *comment
		} else {
			response["approval_comment"] = nil
		}

		result = append(result, response)
	}

	// Final rows error check (IMPORTANT)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func GetMyExpenseRequests(userID int64) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, amount, category, status, reason, approval_comment, created_at
		 FROM expense_requests
		 WHERE employee_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var (
			id        int64
			amount    float64
			category  string
			status    string
			reason    string
			comment   *string
			createdAt time.Time
		)

		if err := rows.Scan(
			&id,
			&amount,
			&category,
			&status,
			&reason,
			&comment,
			&createdAt,
		); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":               id,
			"amount":           amount,
			"category":         category,
			"status":           status,
			"reason":           reason,
			"approval_comment": comment,
			"created_at":       createdAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
func GetMyDiscountRequests(userID int64) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, discount_percentage, status, reason, approval_comment, created_at
		 FROM discount_requests
		 WHERE employee_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	fmt.Print(rows)
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var (
			id      int64
			percent float64
			status  string
			reason  string

			comment *string

			createdAt time.Time
		)

		if err := rows.Scan(
			&id,
			&percent,
			&status,
			&reason,
			&comment,
			&createdAt,
		); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":                  id,
			"discount_percentage": percent,
			"status":              status,
			"reason":              reason,
			"approval_comment":    comment,
			"created_at":          createdAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
