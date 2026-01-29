package helpers

import (
	"errors"

	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5/pgconn"
)

func MapPgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return apperrors.ErrForeignKeyViolation
		case "23505":
			return apperrors.ErrDuplicateEntry
		case "23514":
			return apperrors.ErrCheckConstraintFailed
		default:
			return apperrors.ErrDatabase
		}
	}
	return err
}
