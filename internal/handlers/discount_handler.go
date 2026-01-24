package handlers

import (
	"net/http"
	"rule-based-approval-engine/internal/apperrors"
	"rule-based-approval-engine/internal/response"
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
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized user",
			nil,
		)
		return
	}

	var req DiscountApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request payload",
			err.Error(),
		)
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

	response.Created(
		c,
		message,
		gin.H{
			"status": "PENDING or AUTO_APPROVED",
		},
	)
}

func handleApplyDiscountError(c *gin.Context, err error) {
	switch err {

	case apperrors.ErrInvalidDiscountPercent:
		response.Error(
			c,
			http.StatusBadRequest,
			"discount percentage must be greater than zero",
			nil,
		)

	case apperrors.ErrDiscountLimitExceeded:
		response.Error(
			c,
			http.StatusBadRequest,
			"discount limit exceeded",
			nil,
		)

	case apperrors.ErrDiscountBalanceMissing:
		response.Error(
			c,
			http.StatusNotFound,
			"discount balance not found",
			nil,
		)

	case apperrors.ErrRuleNotFound:
		response.Error(
			c,
			http.StatusInternalServerError,
			"discount approval rules not configured",
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
			"failed to apply discount",
			err.Error(),
		)
	}
}

func CancelDiscount(c *gin.Context) {
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

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid discount request id",
			nil,
		)
		return
	}

	err = services.CancelDiscount(userID, requestID)
	if err != nil {
		handleCancelDiscountError(c, err)
		return
	}

	response.Success(
		c,
		"discount request cancelled successfully",
		nil,
	)
}

func handleCancelDiscountError(c *gin.Context, err error) {
	switch err {

	case apperrors.ErrDiscountRequestNotFound:
		response.Error(
			c,
			http.StatusNotFound,
			"discount request not found",
			nil,
		)

	case apperrors.ErrDiscountCannotCancel:
		response.Error(
			c,
			http.StatusBadRequest,
			"this discount request cannot be cancelled",
			nil,
		)

	default:
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to cancel discount request",
			err.Error(),
		)
	}
}
