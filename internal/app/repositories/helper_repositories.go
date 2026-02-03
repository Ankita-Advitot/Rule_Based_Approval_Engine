package repositories

import (
	"context"
	"time"

	"rule-based-approval-engine/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	helperQueryGetMyLeaves = `SELECT id, leave_type, from_date, to_date, status, reason, approval_comment, created_at
		 FROM leave_requests
		 WHERE employee_id = $1
		 ORDER BY created_at DESC`
	helperQueryGetMyExpenses = `SELECT id, amount, category, status, reason, approval_comment, created_at
		 FROM expense_requests
		 WHERE employee_id = $1
		 ORDER BY created_at DESC`
	helperQueryGetMyDiscounts = `SELECT id, discount_percentage, status, reason, approval_comment, created_at
		 FROM discount_requests
		 WHERE employee_id = $1
		 ORDER BY created_at DESC`
	helperQueryAddHoliday = `INSERT INTO holidays (holiday_date, description, created_by)
		 VALUES ($1,$2,$3)`
	helperQueryGetHolidays   = `SELECT id, holiday_date, description FROM holidays ORDER BY holiday_date`
	helperQueryDeleteHoliday = `DELETE FROM holidays WHERE id=$1`
	helperQueryGetStatusDist = `
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
	`
	helperQueryGetTypeReport = `
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
	`
	helperQueryIsHoliday = `SELECT COUNT(*) FROM holidays WHERE holiday_date=$1`
)

type myRequestsRepository struct {
	db *pgxpool.Pool
}

// NewMyRequestsRepository creates a new instance
func NewMyRequestsRepository(db *pgxpool.Pool) MyRequestsRepository {
	return &myRequestsRepository{db: db}
}

func (r *myRequestsRepository) GetMyLeaveRequests(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(
		ctx,
		helperQueryGetMyLeaves,
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *myRequestsRepository) GetMyExpenseRequests(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(
		ctx,
		helperQueryGetMyExpenses,
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

func (r *myRequestsRepository) GetMyDiscountRequests(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(
		ctx,
		helperQueryGetMyDiscounts,
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
			reason    string
			comment   *string
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

// HolidayRepository implementation is below

type holidayRepository struct {
	db *pgxpool.Pool
}

func (r *holidayRepository) AddHoliday(ctx context.Context, date time.Time, desc string, adminID int64) error {
	_, err := r.db.Exec(
		ctx,
		helperQueryAddHoliday,
		date, desc, adminID,
	)
	return err
}

func (r *holidayRepository) GetHolidays(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(
		ctx,
		helperQueryGetHolidays,
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

		if err := rows.Scan(&id, &d, &desc); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":          id,
			"date":        d.Format("2006-01-02"),
			"description": desc,
		})
	}

	return result, nil
}

func (r *holidayRepository) DeleteHoliday(ctx context.Context, holidayID int64) error {
	_, err := r.db.Exec(
		ctx,
		helperQueryDeleteHoliday,
		holidayID,
	)
	return err
}

type reportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetRequestStatusDistribution(ctx context.Context) (map[string]int, error) {
	rows, err := r.db.Query(ctx, helperQueryGetStatusDist)
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

	return result, rows.Err()
}

func (r *reportRepository) GetRequestsByTypeReport(ctx context.Context) ([]models.RequestTypeReport, error) {
	rows, err := r.db.Query(ctx, helperQueryGetTypeReport)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.RequestTypeReport

	for rows.Next() {
		var r models.RequestTypeReport
		if err := rows.Scan(&r.Type, &r.TotalRequests, &r.AutoApproved); err != nil {
			return nil, err
		}

		if r.TotalRequests > 0 {
			r.AutoApprovedPercent = (float64(r.AutoApproved) / float64(r.TotalRequests)) * 100
		}
		reports = append(reports, r)
	}

	return reports, rows.Err()
}

func (r *holidayRepository) IsHoliday(ctx context.Context, date time.Time) (bool, error) {
	var count int
	err := r.db.QueryRow(
		ctx,
		helperQueryIsHoliday,
		date.Format("2006-01-02"),
	).Scan(&count)
	return count > 0, err
}

func NewHolidayRepository(db *pgxpool.Pool) HolidayRepository {
	return &holidayRepository{db: db}
}
