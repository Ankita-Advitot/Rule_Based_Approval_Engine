package services

import (
	"context"
	"rule-based-approval-engine/internal/database"
)

func GetRequestStatusDistribution() (map[string]int, error) {
	ctx := context.Background()

	rows, err := database.DB.Query(ctx, `
				SELECT status_text, COUNT(*) FROM (
			SELECT status::text AS status_text FROM leave_requests
		) l
		GROUP BY status_text

		UNION ALL

		SELECT status_text, COUNT(*) FROM (
			SELECT status::text AS status_text FROM expense_requests
		) e
		GROUP BY status_text

		UNION ALL

		SELECT status_text, COUNT(*) FROM (
			SELECT status::text AS status_text FROM discount_requests
		) d
		GROUP BY status_text;		
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{
		"approved":      0,
		"rejected":      0,
		"pending":       0,
		"cancelled":     0,
		"auto_rejected": 0,
	}

	for rows.Next() {
		var status string
		var count int

		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}

		switch status {
		case "APPROVED", "AUTO_APPROVED":
			result["approved"] += count
		case "REJECTED":
			result["rejected"] += count
		case "AUTO_REJECTED":
			result["auto_rejected"] += count
		case "PENDING":
			result["pending"] += count
		case "CANCELLED":
			result["cancelled"] += count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type RequestTypeReport struct {
	Type                string  `json:"type"`
	TotalRequests       int     `json:"total_requests"`
	AutoApproved        int     `json:"auto_approved"`
	AutoApprovedPercent float64 `json:"auto_approved_percentage"`
}

func GetRequestsByTypeReport() ([]RequestTypeReport, error) {
	ctx := context.Background()

	rows, err := database.DB.Query(ctx, `
					SELECT 'LEAVE', COUNT(*),
		COUNT(*) FILTER (WHERE status_text='AUTO_APPROVED')
	FROM (SELECT status::text AS status_text FROM leave_requests) l

	UNION ALL

	SELECT 'EXPENSE', COUNT(*),
		COUNT(*) FILTER (WHERE status_text='AUTO_APPROVED')
	FROM (SELECT status::text AS status_text FROM expense_requests) e

	UNION ALL

	SELECT 'DISCOUNT', COUNT(*),
		COUNT(*) FILTER (WHERE status_text='AUTO_APPROVED')
	FROM (SELECT status::text AS status_text FROM discount_requests) d


	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []RequestTypeReport

	for rows.Next() {
		var r RequestTypeReport
		if err := rows.Scan(&r.Type, &r.TotalRequests, &r.AutoApproved); err != nil {
			return nil, err
		}

		if r.TotalRequests > 0 {
			r.AutoApprovedPercent = (float64(r.AutoApproved) / float64(r.TotalRequests)) * 100
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}
