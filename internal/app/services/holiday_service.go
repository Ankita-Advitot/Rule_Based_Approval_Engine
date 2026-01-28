package services

import (
	"context"
	"errors"
	"time"

	"rule-based-approval-engine/internal/database"
)

func ensureAdmin(role string) error {
	if role != "ADMIN" {
		return errors.New("only admin can manage holidays")
	}
	return nil
}

func AddHoliday(role string, adminID int64, date time.Time, desc string) error {
	if err := ensureAdmin(role); err != nil {
		return err
	}

	_, err := database.DB.Exec(
		context.Background(),
		`INSERT INTO holidays (holiday_date, description, created_by)
		 VALUES ($1,$2,$3)`,
		date, desc, adminID,
	)

	return err
}

func GetHolidays(role string) ([]map[string]interface{}, error) {
	if err := ensureAdmin(role); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, holiday_date, description FROM holidays ORDER BY holiday_date`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int64
		var d time.Time
		var desc string

		rows.Scan(&id, &d, &desc)

		result = append(result, map[string]interface{}{
			"id":          id,
			"date":        d.Format("2006-01-02"),
			"description": desc,
		})
	}

	return result, nil
}

func DeleteHoliday(role string, holidayID int64) error {
	if err := ensureAdmin(role); err != nil {
		return err
	}

	_, err := database.DB.Exec(
		context.Background(),
		`DELETE FROM holidays WHERE id=$1`,
		holidayID,
	)

	return err
}
