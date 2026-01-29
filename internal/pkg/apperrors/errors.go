package apperrors

import "errors"

var (
	ErrLeaveBalanceExceeded = errors.New("leave balance exceeded")
	ErrUserNotFound         = errors.New("user not found")
	ErrLeaveBalanceMissing  = errors.New("leave balance not found")
	ErrRuleNotFound         = errors.New("approval rule not configured")
	ErrInvalidLeaveDays     = errors.New("invalid leave days")
	ErrLeaveOverlap         = errors.New(
		"you already have a leave request for this date. first cancel the previous request then only you are allowed to apply",
	)
)
var (
	ErrExpenseBalanceMissing  = errors.New("expense balance not found")
	ErrExpenseLimitExceeded   = errors.New("expense limit exceeded")
	ErrInvalidExpenseAmount   = errors.New("invalid expense amount")
	ErrInvalidExpenseCategory = errors.New("invalid expense category")
)
var (
	ErrInvalidDiscountPercent = errors.New("invalid discount percentage")
	ErrDiscountLimitExceeded  = errors.New("discount limit exceeded")
	ErrDiscountBalanceMissing = errors.New("discount balance not found")
)
var (
	ErrDiscountRequestNotFound = errors.New("discount request not found")
	ErrDiscountCannotCancel    = errors.New("cannot cancel finalized discount request")
)
var (
	ErrUnauthorizedApprover      = errors.New("unauthorized approver")
	ErrDiscountRequestNotPending = errors.New("discount request is not pending")
)
var (
	ErrForeignKeyViolation   = errors.New("foreign key violation")
	ErrDuplicateEntry        = errors.New("duplicate entry")
	ErrCheckConstraintFailed = errors.New("check constraint failed")
	ErrDatabase              = errors.New("database error")

	// --- Transaction errors ---
	ErrTransactionFailed    = errors.New("transaction failed")
	ErrRetryableTransaction = errors.New("retryable transaction error")

	// --- State / consistency errors ---
	ErrNothingToUpdate      = errors.New("nothing to update")
	ErrLeaveRequestNotFound = errors.New("leave request not found")

	// --- Runtime safety ---
	ErrRuleEvaluationFailed = errors.New("rule evaluation failed")
	ErrRequestCancelled     = errors.New("request cancelled")
)
