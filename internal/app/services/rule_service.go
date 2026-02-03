package services

import (
	"context"
	"fmt"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/apperrors"
)

// RuleService handles business logic for rule management
type RuleService struct {
	ruleRepo repositories.RuleRepository
}

// NewRuleService creates a new instance of RuleService
func NewRuleService(ruleRepo repositories.RuleRepository) *RuleService {
	return &RuleService{
		ruleRepo: ruleRepo,
	}
}

// GetRule retrieves a rule by request type and grade ID
func (s *RuleService) GetRule(ctx context.Context, requestType string, gradeID int64) (*models.Rule, error) {
	return s.ruleRepo.GetByTypeAndGrade(ctx, requestType, gradeID)
}

// CreateRule creates or updates a rule (admin only)
func (s *RuleService) CreateRule(ctx context.Context, role string, rule models.Rule) error {
	if role != constants.RoleAdmin {
		return apperrors.ErrUnauthorized
	}

	if rule.RequestType == "" {
		return apperrors.ErrRequestTypeRequired
	}

	if rule.Action == "" {
		return apperrors.ErrActionRequired
	}

	if rule.GradeID == 0 {
		return apperrors.ErrGradeIDRequired
	}

	if rule.Condition == nil || len(rule.Condition) == 0 {
		return apperrors.ErrConditionRequired
	}

	err := s.ruleRepo.Create(ctx, &rule)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}

	return nil
}

// GetRules retrieves all rules (admin only)
func (s *RuleService) GetRules(ctx context.Context, role string) ([]models.Rule, error) {
	if role != constants.RoleAdmin {
		return nil, apperrors.ErrUnauthorized
	}

	return s.ruleRepo.GetAll(ctx)
}

// updates an existing rule (admin only)
func (s *RuleService) UpdateRule(ctx context.Context, role string, ruleID int64, rule models.Rule) error {
	if role != constants.RoleAdmin {
		return apperrors.ErrUnauthorized
	}

	return s.ruleRepo.Update(ctx, ruleID, &rule)
}

// deletes a rule by ID (admin only)
func (s *RuleService) DeleteRule(ctx context.Context, role string, ruleID int64) error {
	if role != constants.RoleAdmin {
		return apperrors.ErrUnauthorized
	}

	return s.ruleRepo.Delete(ctx, ruleID)
}
