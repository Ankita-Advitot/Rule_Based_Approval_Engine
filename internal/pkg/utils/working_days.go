package utils

import (
	"context"
	"time"

	"rule-based-approval-engine/internal/database"
)

// Check if a date is a holiday
func IsHoliday(date time.Time) bool {
	var count int
	_ = database.DB.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM holidays WHERE holiday_date=$1`,
		date.Format("2006-01-02"),
	).Scan(&count)

	return count > 0
}

// Count working days between two dates
func CountWorkingDays(from, to time.Time) int {
	days := 0

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		if IsHoliday(d) {
			continue
		}
		days++
	}

	return days
}
