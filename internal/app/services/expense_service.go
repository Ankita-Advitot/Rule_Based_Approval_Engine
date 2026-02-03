package services

import (
	"context"
	"strings"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// handles business logic for expense requests
type ExpenseService struct {
	expenseReqRepo repositories.ExpenseRequestRepository
	balanceRepo    repositories.BalanceRepository
	ruleService    *RuleService
	userRepo       repositories.UserRepository
	db             *pgxpool.Pool
}

// creates a new instance of ExpenseService
func NewExpenseService(
	expenseReqRepo repositories.ExpenseRequestRepository,
	balanceRepo repositories.BalanceRepository,
	ruleService *RuleService,
	userRepo repositories.UserRepository,
	db *pgxpool.Pool,
) *ExpenseService {
	return &ExpenseService{
		expenseReqRepo: expenseReqRepo,
		balanceRepo:    balanceRepo,
		ruleService:    ruleService,
		userRepo:       userRepo,
		db:             db,
	}
}

// processes an expense application
func (s *ExpenseService) ApplyExpense(
	ctx context.Context,
	userID int64,
	amount float64,
	category string,
	reason string,
) (string, string, error) {
	// validations
	if userID <= 0 {
		return "", "", apperrors.ErrInvalidUser
	}

	if amount <= 0 {
		return "", "", apperrors.ErrInvalidExpenseAmount
	}

	if strings.TrimSpace(category) == "" {
		return "", "", apperrors.ErrInvalidExpenseCategory
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", "", apperrors.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	// expense balance
	remaining, err := s.balanceRepo.GetExpenseBalance(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	if amount > remaining {
		return "", "", apperrors.ErrExpenseLimitExceeded
	}

	// user grade
	gradeID, err := s.userRepo.GetGrade(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	// fetch rule
	rule, err := s.ruleService.GetRule(ctx, "EXPENSE", gradeID)
	if err != nil {
		return "", "", apperrors.ErrRuleNotFound
	}

	// apply rule
	result := helpers.MakeDecision("EXPENSE", rule.Condition, amount)
	status := result.Status
	message := result.Message

	// create request
	expenseReq := &models.ExpenseRequest{
		EmployeeID: userID,
		Amount:     amount,
		Category:   category,
		Reason:     reason,
		Status:     status,
		RuleID:     &rule.ID,
	}

	err = s.expenseReqRepo.Create(ctx, tx, expenseReq)
	if err != nil {
		return "", "", apperrors.ErrInsertFailed
	}

	// deduct if auto-approved
	if status == constants.StatusAutoApproved {
		err = s.balanceRepo.DeductExpenseBalance(ctx, tx, userID, amount)
		if err != nil {
			return "", "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", apperrors.ErrTransactionCommit
	}

	return message, status, nil
}

// cancels an expense request
func (s *ExpenseService) CancelExpense(ctx context.Context, userID, requestID int64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	expenseReq, err := s.expenseReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return err
	}

	// Verify ownership
	if expenseReq.EmployeeID != userID {
		return apperrors.ErrExpenseRequestNotFound
	}

	// reuse CanCancel from apply_cancel_rules.go
	if err := helpers.CanCancel(expenseReq.Status); err != nil {
		return err
	}

	err = s.expenseReqRepo.Cancel(ctx, tx, requestID)
	if err != nil {
		return err
	}

	if expenseReq.Status == constants.StatusAutoApproved {
		err = s.balanceRepo.RestoreExpenseBalance(ctx, tx, userID, expenseReq.Amount)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
