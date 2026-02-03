package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	gradeQueryGetLimits = `SELECT annual_leave_limit, annual_expense_limit, discount_limit_percent
		 FROM grades WHERE id=$1`
)

type gradeRepository struct {
	db *pgxpool.Pool
}

// creates a new instance of GradeRepository
func NewGradeRepository(db *pgxpool.Pool) GradeRepository {
	return &gradeRepository{db: db}
}

func (r *gradeRepository) GetLimits(ctx context.Context, tx pgx.Tx, gradeID int64) (leaveLimit int, expenseLimit float64, discountLimit float64, err error) {
	err = tx.QueryRow(
		ctx,
		gradeQueryGetLimits,
		gradeID,
	).Scan(&leaveLimit, &expenseLimit, &discountLimit)

	return
}
