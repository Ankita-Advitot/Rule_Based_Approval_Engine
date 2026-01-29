package apperrors

import "errors"

// --- Leave-related errors ---
var (
	ErrLeaveBalanceExceeded = errors.New("leave balance exceeded")
	ErrUserNotFound         = errors.New("user not found")
	ErrLeaveBalanceMissing  = errors.New("leave balance not found")
	ErrRuleNotFound         = errors.New("approval rule not configured")
	ErrInvalidLeaveDays     = errors.New("invalid leave days")
	ErrLeaveOverlap         = errors.New(
		"you already have a leave request for this date. first cancel the previous request then only you are allowed to apply",
	)
	ErrLeaveRequestNotFound    = errors.New("leave request not found")
	ErrLeaveVerificationFailed = errors.New("unable to verify existing leave requests")
	ErrLeaveCannotCancel       = errors.New("cannot cancel finalized leave request")
)

// --- Expense-related errors ---
var (
	ErrExpenseBalanceMissing  = errors.New("expense balance not found")
	ErrExpenseLimitExceeded   = errors.New("expense limit exceeded")
	ErrInvalidExpenseAmount   = errors.New("invalid expense amount")
	ErrInvalidExpenseCategory = errors.New("invalid expense category")
	ErrExpenseRequestNotFound = errors.New("expense request not found")
	ErrExpenseCannotCancel    = errors.New("cannot cancel finalized expense request")
)

// --- Discount-related errors ---
var (
	ErrInvalidDiscountPercent  = errors.New("invalid discount percentage")
	ErrDiscountLimitExceeded   = errors.New("discount limit exceeded")
	ErrDiscountBalanceMissing  = errors.New("discount balance not found")
	ErrDiscountRequestNotFound = errors.New("discount request not found")
	ErrDiscountCannotCancel    = errors.New("cannot cancel finalized discount request")
)

// --- Authorization & Approval errors ---
var (
	ErrUnauthorizedApprover      = errors.New("unauthorized approver")
	ErrDiscountRequestNotPending = errors.New("discount request is not pending")
	ErrSelfApprovalNotAllowed    = errors.New("self approval is not allowed")
	ErrUnauthorizedRole          = errors.New("unauthorized role")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrAdminOnly                 = errors.New("only admin can manage this resource")
	ErrRequestNotPending         = errors.New("request not pending")
	ErrEmployeeCannotApprove     = errors.New("employees are not allowed to approve requests")
	ErrManagerNeedsAdmin         = errors.New("managers can only be approved by admin")
	ErrAdminRequestNotAllowed    = errors.New("admin requests are not allowed")
	ErrUnauthorizedApproval      = errors.New("unauthorized approval attempt")
	ErrRequestCannotCancel       = errors.New("cannot cancel finalized request")
)

// --- Authentication errors ---
var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrEmailRequired          = errors.New("email is required")
	ErrPasswordRequired       = errors.New("password is required")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrPasswordHashFailed     = errors.New("password hashing failed")
)

// --- JWT & Token errors ---
var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

// --- Validation errors ---
var (
	ErrInvalidUser      = errors.New("invalid user")
	ErrInvalidDateRange = errors.New("from date cannot be after to date")
)

// --- Database errors ---
var (
	ErrForeignKeyViolation   = errors.New("foreign key violation")
	ErrDuplicateEntry        = errors.New("duplicate entry")
	ErrCheckConstraintFailed = errors.New("check constraint failed")
	ErrDatabase              = errors.New("database error")
	ErrQueryFailed           = errors.New("database query failed")
	ErrUpdateFailed          = errors.New("database update failed")
	ErrInsertFailed          = errors.New("database insert failed")
	ErrDeleteFailed          = errors.New("database delete failed")
)

// --- Transaction errors ---
var (
	ErrTransactionFailed    = errors.New("transaction failed")
	ErrRetryableTransaction = errors.New("retryable transaction error")
	ErrTransactionBegin     = errors.New("unable to start transaction")
	ErrTransactionCommit    = errors.New("failed to commit transaction")
)

// --- Balance & Fetch errors ---
var (
	ErrBalanceUpdateFailed = errors.New("failed to update balance")
	ErrBalanceFetchFailed  = errors.New("failed to fetch balance")
)

// --- State / consistency errors ---
var (
	ErrNothingToUpdate = errors.New("nothing to update")
)

// --- Runtime safety ---
var (
	ErrRuleEvaluationFailed = errors.New("rule evaluation failed")
	ErrRequestCancelled     = errors.New("request cancelled")
)

// --- Rule Service errors ---
var (
	ErrNoRuleFound           = errors.New("no rule found")
	ErrRequestTypeRequired   = errors.New("request_type is required")
	ErrActionRequired        = errors.New("action is required")
	ErrGradeIDRequired       = errors.New("grade_id is required")
	ErrConditionRequired     = errors.New("condition is required")
	ErrInvalidConditionJSON  = errors.New("invalid condition JSON")
	ErrRuleNotFoundForDelete = errors.New("rule not found")
)
