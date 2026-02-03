package services

import (
	"context"
	"time"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// handles business logic for leave requests
type LeaveService struct {
	leaveReqRepo repositories.LeaveRequestRepository
	balanceRepo  repositories.BalanceRepository
	ruleService  *RuleService
	userRepo     repositories.UserRepository
	db           *pgxpool.Pool
}

// creates a new instance of LeaveService
func NewLeaveService(
	leaveReqRepo repositories.LeaveRequestRepository,
	balanceRepo repositories.BalanceRepository,
	ruleService *RuleService,
	userRepo repositories.UserRepository,
	db *pgxpool.Pool,
) *LeaveService {
	return &LeaveService{
		leaveReqRepo: leaveReqRepo,
		balanceRepo:  balanceRepo,
		ruleService:  ruleService,
		userRepo:     userRepo,
		db:           db,
	}
}

// processes a leave application
func (s *LeaveService) ApplyLeave(
	ctx context.Context,
	userID int64,
	from time.Time,
	to time.Time,
	days int,
	leaveType string,
	reason string,
) (string, string, error) {
	// validations
	if userID <= 0 {
		return "", "", apperrors.ErrInvalidUser
	}

	if days <= 0 {
		return "", "", apperrors.ErrInvalidLeaveDays
	}

	if from.After(to) {
		return "", "", apperrors.ErrInvalidDateRange
	}

	// date validation
	today := time.Now().Truncate(24 * time.Hour)
	if from.Before(today) {
		return "", "", apperrors.ErrPastDate
	}

	// check overlap
	overlap, err := s.leaveReqRepo.CheckOverlap(ctx, userID, from, to)
	if err != nil {
		return "", "", apperrors.ErrLeaveVerificationFailed
	}

	if overlap {
		return "", "", apperrors.ErrLeaveOverlap
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", "", apperrors.ErrTransactionBegin
	}
	defer tx.Rollback(ctx)

	// leave balance
	remaining, err := s.balanceRepo.GetLeaveBalance(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	if days > remaining {
		return "", "", apperrors.ErrLeaveBalanceExceeded
	}

	// user grade
	gradeID, err := s.userRepo.GetGrade(ctx, tx, userID)
	if err != nil {
		return "", "", err
	}

	// fetch rule
	rule, err := s.ruleService.GetRule(ctx, "LEAVE", gradeID)
	if err != nil {
		return "", "", apperrors.ErrRuleNotFound
	}

	// apply rule
	result := helpers.MakeDecision("LEAVE", rule.Condition, float64(days))
	status := result.Status
	message := result.Message

	leaveReq := &models.LeaveRequest{
		EmployeeID: userID,
		FromDate:   from,
		ToDate:     to,
		Reason:     reason,
		LeaveType:  leaveType,
		Status:     status,
		RuleID:     &rule.ID,
	}

	err = s.leaveReqRepo.Create(ctx, tx, leaveReq)
	if err != nil {
		return "", "", helpers.MapPgError(err)
	}

	// deduct if auto-approved
	if status == constants.StatusApproved {
		err = s.balanceRepo.DeductLeaveBalance(ctx, tx, userID, days)
		if err != nil {
			return "", "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", apperrors.ErrTransactionCommit
	}

	return message, status, nil
}

// cancels a leave request
func (s *LeaveService) CancelLeave(ctx context.Context, userID, requestID int64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	leaveReq, err := s.leaveReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return err
	}

	// Verify ownership
	if leaveReq.EmployeeID != userID {
		return apperrors.ErrLeaveRequestNotFound
	}

	// reuse CanCancel from apply_cancel_rules.go
	if err := helpers.CanCancel(leaveReq.Status); err != nil {
		return err
	}

	days := utils.CalculateLeaveDays(leaveReq.FromDate, leaveReq.ToDate)

	err = s.leaveReqRepo.Cancel(ctx, tx, requestID)
	if err != nil {
		return err
	}

	if leaveReq.Status == constants.StatusAutoApproved {
		err = s.balanceRepo.RestoreLeaveBalance(ctx, tx, userID, days)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
