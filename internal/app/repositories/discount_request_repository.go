package repositories

import (
	"context"
	"rule-based-approval-engine/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	discountQueryCreate = `INSERT INTO discount_requests
		 (employee_id, discount_percentage, reason, status, rule_id)
		 VALUES ($1, $2, $3, $4, $5)`
	discountQueryGetByID = `SELECT id, employee_id, discount_percentage, reason, status, rule_id, approved_by_id, created_at
		 FROM discount_requests WHERE id=$1`
	discountQueryUpdateStatus = `UPDATE discount_requests
		 SET status=$1, approved_by_id=$2, approval_comment=$3
		 WHERE id=$4`
	discountQueryGetPendingForManager = `
		SELECT dr.id, u.name, dr.discount_percentage, dr.reason, dr.created_at
		FROM discount_requests dr
		JOIN users u ON dr.employee_id = u.id
		WHERE dr.status='PENDING' AND u.manager_id=$1
	`
	discountQueryGetPendingForAdmin = `
		SELECT dr.id, u.name, dr.discount_percentage, dr.reason, dr.created_at
		FROM discount_requests dr
		JOIN users u ON dr.employee_id = u.id
		WHERE dr.status='PENDING'
	`
	discountQueryCancel             = `UPDATE discount_requests SET status='CANCELLED' WHERE id=$1`
	discountQueryGetPendingRequests = "SELECT id, created_at FROM discount_requests WHERE status='PENDING'"
)

type discountRequestRepository struct {
	db *pgxpool.Pool
}

func NewDiscountRequestRepository(db *pgxpool.Pool) DiscountRequestRepository {
	return &discountRequestRepository{db: db}
}

func (r *discountRequestRepository) Create(ctx context.Context, tx pgx.Tx, req *models.DiscountRequest) error {
	_, err := tx.Exec(
		ctx,
		discountQueryCreate,
		req.EmployeeID, req.DiscountPercentage, req.Reason, req.Status, req.RuleID,
	)
	return err
}

func (r *discountRequestRepository) GetByID(ctx context.Context, tx pgx.Tx, requestID int64) (*models.DiscountRequest, error) {
	req := &models.DiscountRequest{}
	err := tx.QueryRow(
		ctx,
		discountQueryGetByID,
		requestID,
	).Scan(&req.ID, &req.EmployeeID, &req.DiscountPercentage, &req.Reason, &req.Status, &req.RuleID, &req.ApprovedByID, &req.CreatedAt)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *discountRequestRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, requestID int64, status string, approverID int64, comment string) error {
	_, err := tx.Exec(
		ctx,
		discountQueryUpdateStatus,
		status, approverID, comment, requestID,
	)
	return err
}

func (r *discountRequestRepository) GetPendingForManager(ctx context.Context, managerID int64) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, discountQueryGetPendingForManager, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int64
		var name, reason string
		var percent float64
		var createdAt interface{}
		if err := rows.Scan(&id, &name, &percent, &reason, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"id":                  id,
			"employee":            name,
			"discount_percentage": percent,
			"reason":              reason,
			"created_at":          createdAt,
		})
	}
	return result, nil
}

func (r *discountRequestRepository) GetPendingForAdmin(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, discountQueryGetPendingForAdmin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int64
		var name, reason string
		var percent float64
		var createdAt interface{}
		if err := rows.Scan(&id, &name, &percent, &reason, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"id":                  id,
			"employee":            name,
			"discount_percentage": percent,
			"reason":              reason,
			"created_at":          createdAt,
		})
	}
	return result, nil
}

func (r *discountRequestRepository) Cancel(ctx context.Context, tx pgx.Tx, requestID int64) error {
	_, err := r.db.Exec(ctx, discountQueryCancel, requestID)
	return err
}

func (r *discountRequestRepository) GetPendingRequests(ctx context.Context) ([]struct {
	ID        int64
	CreatedAt time.Time
}, error) {
	rows, err := r.db.Query(ctx, discountQueryGetPendingRequests)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []struct {
		ID        int64
		CreatedAt time.Time
	}
	for rows.Next() {
		var item struct {
			ID        int64
			CreatedAt time.Time
		}
		if err := rows.Scan(&item.ID, &item.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}
