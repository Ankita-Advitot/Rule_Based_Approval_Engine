package handlers

import (
	"net/http"
	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExpenseApplyRequest struct {
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Reason   string  `json:"reason"`
}

func ApplyExpense(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized user",
			nil,
		)
		return
	}

	var req ExpenseApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request payload",
			err.Error(),
		)
		return
	}

	message, status, err := services.ApplyExpense(
		userID,
		req.Amount,
		req.Category,
		req.Reason,
	)

	if err != nil {
		handleApplyExpenseError(c, err)
		return
	}

	response.Created(
		c,
		message,
		gin.H{
			"status": status,
		},
	)
}
func handleApplyExpenseError(c *gin.Context, err error) {
	switch err {

	case apperrors.ErrInvalidExpenseAmount:
		response.Error(
			c,
			http.StatusBadRequest,
			"expense amount must be greater than zero",
			nil,
		)

	case apperrors.ErrInvalidExpenseCategory:
		response.Error(
			c,
			http.StatusBadRequest,
			"expense category is required",
			nil,
		)

	case apperrors.ErrExpenseLimitExceeded:
		response.Error(
			c,
			http.StatusBadRequest,
			"expense limit exceeded",
			nil,
		)

	case apperrors.ErrExpenseBalanceMissing:
		response.Error(
			c,
			http.StatusNotFound,
			"expense balance not found",
			nil,
		)

	case apperrors.ErrRuleNotFound:
		response.Error(
			c,
			http.StatusInternalServerError,
			"expense approval rules not configured",
			nil,
		)

	case apperrors.ErrUserNotFound:
		response.Error(
			c,
			http.StatusNotFound,
			"user not found",
			nil,
		)

	default:
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to apply expense",
			err.Error(),
		)
	}
}

func CancelExpense(c *gin.Context) {
	userID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid expense request id",
			nil,
		)
		return
	}

	err = services.CancelExpense(userID, requestID)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"unable to cancel expense request",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"expense request cancelled successfully",
		nil,
	)
}
