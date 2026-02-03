package services

import (
	"context"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/app/services/helpers"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DiscountApprovalService struct {
	discountReqRepo repositories.DiscountRequestRepository
	balanceRepo     repositories.BalanceRepository
	userRepo        repositories.UserRepository
	db              *pgxpool.Pool
}

func NewDiscountApprovalService(
	discountReqRepo repositories.DiscountRequestRepository,
	balanceRepo repositories.BalanceRepository,
	userRepo repositories.UserRepository,
	db *pgxpool.Pool,
) *DiscountApprovalService {
	return &DiscountApprovalService{
		discountReqRepo: discountReqRepo,
		balanceRepo:     balanceRepo,
		userRepo:        userRepo,
		db:              db,
	}
}

func (s *DiscountApprovalService) GetPendingDiscountRequests(ctx context.Context, role string, approverID int64) ([]map[string]interface{}, error) {
	if role == constants.RoleManager {
		return s.discountReqRepo.GetPendingForManager(ctx, approverID)
	} else if role == constants.RoleAdmin {
		return s.discountReqRepo.GetPendingForAdmin(ctx)
	} else {
		return nil, apperrors.ErrUnauthorized
	}
}

func (s *DiscountApprovalService) ApproveDiscount(ctx context.Context, role string, approverID, requestID int64, comment string) error {
	if role == constants.RoleEmployee {
		return apperrors.ErrEmployeeCannotApprove
	}

	if comment == "" {
		return apperrors.ErrCommentRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	discountReq, err := s.discountReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return apperrors.ErrDiscountRequestNotFound
	}

	if approverID == discountReq.EmployeeID {
		return apperrors.ErrSelfApprovalNotAllowed
	}

	if err := helpers.ValidatePendingStatus(discountReq.Status); err != nil {
		return err
	}

	requesterRole, err := s.userRepo.GetRole(ctx, tx, discountReq.EmployeeID)
	if err != nil {
		return err
	}

	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// Update request
	err = s.discountReqRepo.UpdateStatus(ctx, tx, requestID, "APPROVED", approverID, comment)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *DiscountApprovalService) RejectDiscount(ctx context.Context, role string, approverID, requestID int64, comment string) error {
	if role == constants.RoleEmployee {
		return apperrors.ErrEmployeeCannotApprove
	}

	if comment == "" {
		return apperrors.ErrCommentRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	discountReq, err := s.discountReqRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return apperrors.ErrDiscountRequestNotFound
	}

	if approverID == discountReq.EmployeeID {
		return apperrors.ErrSelfApprovalNotAllowed
	}

	if err := helpers.ValidatePendingStatus(discountReq.Status); err != nil {
		return err
	}

	requesterRole, err := s.userRepo.GetRole(ctx, tx, discountReq.EmployeeID)
	if err != nil {
		return err
	}

	if err := helpers.ValidateApproverRole(role, requesterRole); err != nil {
		return err
	}

	// Update request
	err = s.discountReqRepo.UpdateStatus(ctx, tx, requestID, "REJECTED", approverID, comment)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
