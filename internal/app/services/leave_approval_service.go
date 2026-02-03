package services

import (
	"context"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// handles business logic for leave approval operations
type LeaveApprovalService struct {
	leaveReqRepo repositories.LeaveRequestRepository
	balanceRepo  repositories.BalanceRepository
	userRepo     repositories.UserRepository
	db           *pgxpool.Pool
}

// creates a new instance of LeaveApprovalService
func NewLeaveApprovalService(
	leaveReqRepo repositories.LeaveRequestRepository,
	balanceRepo repositories.BalanceRepository,
	userRepo repositories.UserRepository,
	db *pgxpool.Pool,
) *LeaveApprovalService {
	return &LeaveApprovalService{
		leaveReqRepo: leaveReqRepo,
		balanceRepo:  balanceRepo,
		userRepo:     userRepo,
		db:           db,
	}
}

// retrieves pending leave requests based on role
func (s *LeaveApprovalService) GetPendingLeaveRequests(ctx context.Context, role string, approverID int64) ([]map[string]interface{}, error) {
	switch role {
	case constants.RoleManager:
		return s.leaveReqRepo.GetPendingForManager(ctx, approverID)
	case constants.RoleAdmin:
		return s.leaveReqRepo.GetPendingForAdmin(ctx)
	default:
		return nil, apperrors.ErrUnauthorizedRole
	}
}

// approves a leave request
func (s *LeaveApprovalService) ApproveLeave(
	ctx context.Context,
	role string,
	approverID, requestID int64,
	approvalComment string,
) error {
	// check role
	if role == constants.RoleEmployee {
		return apperrors.ErrEmployeeCannotApprove
	}

	// validate comment
	if approvalComment == "" {
		return apperrors.ErrCommentRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	leaveReq, err := s.leaveReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return err
	}

	if approverID == leaveReq.EmployeeID {
		return apperrors.ErrSelfApprovalNotAllowed
	}

	if err := helpers.ValidatePendingStatus(leaveReq.Status); err != nil {
		return err
	}

	// Authorization against requester
	requesterRole, err := s.userRepo.GetRole(ctx, tx, leaveReq.EmployeeID)
	if err != nil {
		return err
	}

	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	days := utils.CalculateLeaveDays(leaveReq.FromDate, leaveReq.ToDate)

	// Deduct leave balance
	err = s.balanceRepo.DeductLeaveBalance(ctx, tx, leaveReq.EmployeeID, days)
	if err != nil {
		return err
	}

	// Default comment if not provided
	if approvalComment == "" {
		approvalComment = constants.StatusApproved
	}

	// Update request
	err = s.leaveReqRepo.UpdateStatus(ctx, tx, requestID, "APPROVED", approverID, approvalComment)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RejectLeave rejects a leave request
func (s *LeaveApprovalService) RejectLeave(
	ctx context.Context,
	role string,
	approverID, requestID int64,
	rejectionComment string,
) error {
	// check role
	if role == constants.RoleEmployee {
		return apperrors.ErrEmployeeCannotApprove
	}

	// validate comment
	if rejectionComment == "" {
		return apperrors.ErrCommentRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	leaveReq, err := s.leaveReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return err
	}

	if err := helpers.ValidatePendingStatus(leaveReq.Status); err != nil {
		return err
	}

	if approverID == leaveReq.EmployeeID {
		return apperrors.ErrSelfApprovalNotAllowed
	}

	// Authorization
	requesterRole, err := s.userRepo.GetRole(ctx, tx, leaveReq.EmployeeID)
	if err != nil {
		return err
	}

	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// Update request
	err = s.leaveReqRepo.UpdateStatus(ctx, tx, requestID, "REJECTED", approverID, rejectionComment)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
