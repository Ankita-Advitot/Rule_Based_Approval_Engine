package handlers

import (
	"net/http"
	"rule-based-approval-engine/internal/services"
	"strconv"
	"rule-based-approval-engine/internal/apperrors"

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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	var req ExpenseApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	message, err := services.ApplyExpense(
		userID,
		req.Amount,
		req.Category,
		req.Reason,
	)

	if err != nil {
		handleApplyExpenseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}
func handleApplyExpenseError(c *gin.Context, err error) {

	switch err {

	case apperrors.ErrInvalidExpenseAmount:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "expense amount must be greater than zero",
		})

	case apperrors.ErrInvalidExpenseCategory:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "expense category is required",
		})

	case apperrors.ErrExpenseLimitExceeded:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "expense limit exceeded",
		})

	case apperrors.ErrExpenseBalanceMissing:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "expense balance not found",
		})

	case apperrors.ErrRuleNotFound:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "expense approval rules not configured",
		})

	case apperrors.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to apply expense",
		})
	}
}

func CancelExpense(c *gin.Context) {
	userID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request id"})
		return
	}

	err = services.CancelExpense(userID, requestID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "expense request cancelled"})
}
