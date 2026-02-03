package services

import (
	"context"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DiscountService struct {
	discountReqRepo repositories.DiscountRequestRepository
	balanceRepo     repositories.BalanceRepository
	ruleService     *RuleService
	userRepo        repositories.UserRepository
	db              *pgxpool.Pool
}

func NewDiscountService(
	discountReqRepo repositories.DiscountRequestRepository,
	balanceRepo repositories.BalanceRepository,
	ruleService *RuleService,
	userRepo repositories.UserRepository,
	db *pgxpool.Pool,
) *DiscountService {
	return &DiscountService{
		discountReqRepo: discountReqRepo,
		balanceRepo:     balanceRepo,
		ruleService:     ruleService,
		userRepo:        userRepo,
		db:              db,
	}
}

func (s *DiscountService) ApplyDiscount(
	ctx context.Context,
	userID int64,
	percent float64,
	reason string,
) (string, string, error) {

	// validations
	if userID <= 0 {
		return "", "", apperrors.ErrInvalidUser
	}

	if percent <= 0 {
		return "", "", apperrors.ErrInvalidDiscountPercent
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", "", apperrors.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	// fetch remaining
	remaining, err := s.balanceRepo.GetDiscountBalance(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	if percent > remaining {
		return "", "", apperrors.ErrDiscountLimitExceeded
	}

	// user grade
	gradeID, err := s.userRepo.GetGrade(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	// fetch rule
	rule, err := s.ruleService.GetRule(ctx, "DISCOUNT", gradeID)
	if err != nil {
		return "", "", apperrors.ErrRuleNotFound
	}

	// apply rule
	result := helpers.MakeDecision("DISCOUNT", rule.Condition, percent)
	status := result.Status
	message := result.Message

	// create request
	discountReq := &models.DiscountRequest{
		EmployeeID:         userID,
		DiscountPercentage: percent,
		Reason:             reason,
		Status:             status,
		RuleID:             &rule.ID,
	}

	err = s.discountReqRepo.Create(ctx, tx, discountReq)
	if err != nil {
		return "", "", apperrors.ErrInsertFailed
	}

	// deduct if auto-approved
	if status == constants.StatusApproved {
		err = s.balanceRepo.DeductDiscountBalance(ctx, tx, userID, percent)
		if err != nil {
			return "", "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", apperrors.ErrTransactionCommit
	}

	return message, status, nil
}

func (s *DiscountService) CancelDiscount(ctx context.Context, userID, requestID int64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return apperrors.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	discountReq, err := s.discountReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return apperrors.ErrDiscountRequestNotFound
	}

	if discountReq.EmployeeID != userID {
		return apperrors.ErrDiscountRequestNotFound
	}

	// reuse CanCancel from apply_cancel_rules.go
	if err := helpers.CanCancel(discountReq.Status); err != nil {
		return err
	}

	err = s.discountReqRepo.Cancel(ctx, tx, requestID)
	if err != nil {
		return apperrors.ErrUpdateFailed
	}

	if discountReq.Status == constants.StatusAutoApproved {
		err = s.balanceRepo.RestoreDiscountBalance(ctx, tx, userID, discountReq.DiscountPercentage)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
