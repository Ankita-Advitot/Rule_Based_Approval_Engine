package services

import (
	"context"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ExpenseApprovalService handles business logic for expense approval operations
type ExpenseApprovalService struct {
	expenseReqRepo repositories.ExpenseRequestRepository
	balanceRepo    repositories.BalanceRepository
	userRepo       repositories.UserRepository
	db             *pgxpool.Pool
}

// NewExpenseApprovalService creates a new instance of ExpenseApprovalService
func NewExpenseApprovalService(
	expenseReqRepo repositories.ExpenseRequestRepository,
	balanceRepo repositories.BalanceRepository,
	userRepo repositories.UserRepository,
	db *pgxpool.Pool,
) *ExpenseApprovalService {
	return &ExpenseApprovalService{
		expenseReqRepo: expenseReqRepo,
		balanceRepo:    balanceRepo,
		userRepo:       userRepo,
		db:             db,
	}
}

// GetPendingExpenseRequests retrieves pending expense requests based on role
func (s *ExpenseApprovalService) GetPendingExpenseRequests(ctx context.Context, role string, approverID int64) ([]map[string]interface{}, error) {
	if role == constants.RoleManager {
		return s.expenseReqRepo.GetPendingForManager(ctx, approverID)
	} else if role == constants.RoleAdmin {
		return s.expenseReqRepo.GetPendingForAdmin(ctx)
	} else {
		return nil, apperrors.ErrUnauthorized
	}
}

// approves an expense request
func (s *ExpenseApprovalService) ApproveExpense(
	ctx context.Context,
	role string,
	approverID, requestID int64,
	comment string,
) error {
	// check role
	if role == constants.RoleEmployee {
		return apperrors.ErrEmployeeCannotApprove
	}

	// validate comment
	if comment == "" {
		return apperrors.ErrCommentRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	expenseReq, err := s.expenseReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return err
	}

	if approverID == expenseReq.EmployeeID {
		return apperrors.ErrSelfApprovalNotAllowed
	}

	// validate status
	if err := helpers.ValidatePendingStatus(expenseReq.Status); err != nil {
		return err
	}

	// fetch role
	requesterRole, err := s.userRepo.GetRole(ctx, tx, expenseReq.EmployeeID)
	if err != nil {
		return err
	}

	// validate role
	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// update request
	err = s.expenseReqRepo.UpdateStatus(ctx, tx, requestID, "APPROVED", approverID, comment)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// rejects an expense request
func (s *ExpenseApprovalService) RejectExpense(
	ctx context.Context,
	role string,
	approverID, requestID int64,
	comment string,
) error {
	// check role
	if role == constants.RoleEmployee {
		return apperrors.ErrEmployeeCannotApprove
	}

	// validate comment
	if comment == "" {
		return apperrors.ErrCommentRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	expenseReq, err := s.expenseReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return err
	}

	if approverID == expenseReq.EmployeeID {
		return apperrors.ErrSelfApprovalNotAllowed
	}

	// 2. Validate pending status
	if err := helpers.ValidatePendingStatus(expenseReq.Status); err != nil {
		return err
	}

	// 4. Fetch requester role
	requesterRole, err := s.userRepo.GetRole(ctx, tx, expenseReq.EmployeeID)
	if err != nil {
		return err
	}

	// 5. Validate approver role
	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// 6. Default rejection comment
	if comment == "" {
		comment = "Rejected"
	}

	// 7. Update request
	err = s.expenseReqRepo.UpdateStatus(ctx, tx, requestID, "REJECTED", approverID, comment)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
