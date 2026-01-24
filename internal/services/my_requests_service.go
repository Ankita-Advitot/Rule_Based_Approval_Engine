package services

import (
	"context"
	"time"

	"rule-based-approval-engine/internal/database"
)

func GetMyLeaveRequests(userID int64) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, leave_type, from_date, to_date, status, approval_comment, created_at
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
			comment   *string
			createdAt time.Time
		)

		if err := rows.Scan(
			&id,
			&leaveType,
			&fromDate,
			&toDate,
			&status,
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
		`SELECT id, amount, category, status, approval_comment, created_at
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
			comment   *string
			createdAt time.Time
		)

		rows.Scan(&id, &amount, &category, &status, &comment, &createdAt)

		result = append(result, map[string]interface{}{
			"id":               id,
			"amount":           amount,
			"category":         category,
			"status":           status,
			"approval_comment": comment,
			"created_at":       createdAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
func GetMyDiscountRequests(userID int64) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, discount_percentage, status, approval_comment, created_at
		 FROM discount_requests
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
			percent   float64
			status    string
			comment   *string
			createdAt string
		)

		rows.Scan(&id, &percent, &status, &comment, &createdAt)

		result = append(result, map[string]interface{}{
			"id":                  id,
			"discount_percentage": percent,
			"status":              status,
			"approval_comment":    comment,
			"created_at":          createdAt,
		})
	}

	return result, nil
}
