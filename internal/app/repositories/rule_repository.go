package repositories

import (
	"context"
	"encoding/json"

	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ruleQueryGetByTypeAndGrade = `SELECT id, condition, action 
		 FROM rules 
		 WHERE request_type=$1 AND grade_id=$2 AND active=true
		 LIMIT 1`
	ruleQueryCreate = `INSERT INTO rules (request_type, condition, action, grade_id, active)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (request_type, grade_id)
		 DO UPDATE SET
		 	condition  = EXCLUDED.condition,
		 	action     = EXCLUDED.action,
		 	active     = EXCLUDED.active,
		 	updated_at = NOW()`
	ruleQueryGetAll = `SELECT id, request_type, condition, action, grade_id, active
		 FROM rules`
	ruleQueryUpdate = `UPDATE rules
		 SET request_type=$1,
		     condition=$2,
		     action=$3,
		     grade_id=$4,
		     active=$5
		 WHERE id=$6`
	ruleQueryDelete = `DELETE FROM rules WHERE id=$1`
)

type ruleRepository struct {
	db *pgxpool.Pool
}

// instance
func NewRuleRepository(db *pgxpool.Pool) RuleRepository {
	return &ruleRepository{db: db}
}

func (r *ruleRepository) GetByTypeAndGrade(ctx context.Context, requestType string, gradeID int64) (*models.Rule, error) {
	var rule models.Rule
	var conditionJSON []byte

	err := r.db.QueryRow(
		ctx,
		ruleQueryGetByTypeAndGrade,
		requestType, gradeID,
	).Scan(&rule.ID, &conditionJSON, &rule.Action)

	if err != nil {
		return nil, apperrors.ErrNoRuleFound
	}

	err = json.Unmarshal(conditionJSON, &rule.Condition)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (r *ruleRepository) Create(ctx context.Context, rule *models.Rule) error {
	conditionJSON, err := json.Marshal(rule.Condition)
	if err != nil {
		return apperrors.ErrInvalidConditionJSON
	}

	_, err = r.db.Exec(
		ctx,
		ruleQueryCreate,
		rule.RequestType,
		conditionJSON,
		rule.Action,
		rule.GradeID,
		rule.Active,
	)

	return err
}

func (r *ruleRepository) GetAll(ctx context.Context) ([]models.Rule, error) {
	rows, err := r.db.Query(
		ctx,
		ruleQueryGetAll,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.Rule

	for rows.Next() {
		var rule models.Rule
		var conditionJSON []byte

		if err := rows.Scan(
			&rule.ID,
			&rule.RequestType,
			&conditionJSON,
			&rule.Action,
			&rule.GradeID,
			&rule.Active,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal(conditionJSON, &rule.Condition)
		rules = append(rules, rule)
	}

	return rules, nil
}

func (r *ruleRepository) Update(ctx context.Context, ruleID int64, rule *models.Rule) error {
	conditionJSON, err := json.Marshal(rule.Condition)
	if err != nil {
		return err
	}

	cmd, err := r.db.Exec(
		ctx,
		ruleQueryUpdate,
		rule.RequestType,
		conditionJSON,
		rule.Action,
		rule.GradeID,
		rule.Active,
		ruleID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return apperrors.ErrNoRuleFound
	}

	return nil
}

func (r *ruleRepository) Delete(ctx context.Context, ruleID int64) error {
	cmd, err := r.db.Exec(
		ctx,
		ruleQueryDelete,
		ruleID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return apperrors.ErrRuleNotFoundForDelete
	}

	return nil
}
