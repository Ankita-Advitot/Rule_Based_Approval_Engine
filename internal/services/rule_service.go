package services

import (
	"context"
	"encoding/json"
	"errors"

	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/models"
)

func GetRule(requestType string, gradeID int64) (*models.Rule, error) {
	var rule models.Rule
	var conditionJSON []byte

	err := database.DB.QueryRow(
		context.Background(),
		`SELECT id, condition, action 
		 FROM rules 
		 WHERE request_type=$1 AND grade_id=$2 AND active=true
		 LIMIT 1`,
		requestType, gradeID,
	).Scan(&rule.ID, &conditionJSON, &rule.Action)

	if err != nil {
		return nil, errors.New("no rule found")
	}

	err = json.Unmarshal(conditionJSON, &rule.Condition)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func CreateRule(role string, rule models.Rule) error {
	if role != "ADMIN" {
		return errors.New("unauthorized")
	}

	ctx := context.Background()

	conditionJSON, err := json.Marshal(rule.Condition)
	if err != nil {
		return err
	}
	if rule.RequestType == "" {
		return errors.New("request_type is required")
	}
	if rule.Action == "" {
		return errors.New("action is required")
	}
	if rule.GradeID == 0 {
		return errors.New("grade_id is required")
	}

	_, err = database.DB.Exec(
		ctx,
		`INSERT INTO rules (request_type, condition, action, grade_id, active)
		 VALUES ($1, $2, $3, $4, $5)`,
		rule.RequestType,
		conditionJSON,
		rule.Action,
		rule.GradeID,
		rule.Active,
	)

	return err
}

func GetRules(role string) ([]models.Rule, error) {
	if role != "ADMIN" {
		return nil, errors.New("unauthorized")
	}

	ctx := context.Background()

	rows, err := database.DB.Query(
		ctx,
		`SELECT id, request_type, condition, action, grade_id, active
		 FROM rules`,
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

func UpdateRule(role string, ruleID int64, rule models.Rule) error {
	if role != "ADMIN" {
		return errors.New("unauthorized")
	}

	ctx := context.Background()

	conditionJSON, err := json.Marshal(rule.Condition)
	if err != nil {
		return err
	}

	cmd, err := database.DB.Exec(
		ctx,
		`UPDATE rules
		 SET request_type=$1,
		     condition=$2,
		     action=$3,
		     grade_id=$4,
		     active=$5
		 WHERE id=$6`,
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
		return errors.New("rule not found")
	}

	return nil
}

func DeleteRule(role string, ruleID int64) error {
	if role != "ADMIN" {
		return errors.New("unauthorized")
	}

	ctx := context.Background()

	cmd, err := database.DB.Exec(
		ctx,
		`DELETE FROM rules WHERE id=$1`,
		ruleID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("rule not found")
	}

	return nil
}
