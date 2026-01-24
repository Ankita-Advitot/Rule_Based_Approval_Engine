package handlers

import (
	"net/http"
	"rule-based-approval-engine/internal/apperrors"
	"rule-based-approval-engine/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DiscountApplyRequest struct {
	DiscountPercentage float64 `json:"discount_percentage"`
	Reason             string  `json:"reason"`
}

func ApplyDiscount(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	var req DiscountApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	message, err := services.ApplyDiscount(
		userID,
		req.DiscountPercentage,
		req.Reason,
	)

	if err != nil {
		handleApplyDiscountError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

func handleApplyDiscountError(c *gin.Context, err error) {

	switch err {

	case apperrors.ErrInvalidDiscountPercent:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "discount percentage must be greater than zero",
		})

	case apperrors.ErrDiscountLimitExceeded:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "discount limit exceeded",
		})

	case apperrors.ErrDiscountBalanceMissing:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "discount balance not found",
		})

	case apperrors.ErrRuleNotFound:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "discount approval rules not configured",
		})

	case apperrors.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to apply discount",
		})
	}
}
func CancelDiscount(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request id",
		})
		return
	}

	err = services.CancelDiscount(userID, requestID)
	if err != nil {
		handleCancelDiscountError(c, err) // ✅ IMPORTANT
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "discount request cancelled",
	})
}

func handleCancelDiscountError(c *gin.Context, err error) {
	switch err {

	case apperrors.ErrDiscountRequestNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "discount request not found",
		})

	case apperrors.ErrDiscountCannotCancel:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "this discount request cannot be cancelled",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to cancel discount request",
		})
	}
}
