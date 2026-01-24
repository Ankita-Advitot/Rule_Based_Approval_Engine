package handlers

import (
	"net/http"
	"strconv"

	"rule-based-approval-engine/internal/apperrors"
	"rule-based-approval-engine/internal/response"
	"rule-based-approval-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func GetPendingDiscounts(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	discounts, err := services.GetPendingDiscountRequests(role, userID)
	if err != nil {
		response.Error(
			c,
			http.StatusForbidden,
			"failed to fetch pending discount requests",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"pending discount requests fetched successfully",
		discounts,
	)
}

func ApproveDiscount(c *gin.Context) {
	handleDiscountDecision(c, "APPROVED")
}

func RejectDiscount(c *gin.Context) {
	handleDiscountDecision(c, "REJECTED")
}

type ApprovalRequest struct {
	Comment string `json:"comment"`
}

func handleDiscountDecision(c *gin.Context, action string) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

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

	var req ApprovalRequest
	_ = c.ShouldBindJSON(&req) // comment is optional

	var svcErr error
	if action == "APPROVED" {
		svcErr = services.ApproveDiscount(role, approverID, requestID, req.Comment)
	} else {
		svcErr = services.RejectDiscount(role, approverID, requestID, req.Comment)
	}

	if svcErr != nil {
		handleApproveRejectDiscountError(c, svcErr)
		return
	}

	response.Success(
		c,
		"discount request "+action+" successfully",
		nil,
	)
}

func handleApproveRejectDiscountError(c *gin.Context, err error) {
	switch err {

	case apperrors.ErrUnauthorizedApprover:
		response.Error(
			c,
			http.StatusForbidden,
			"you are not authorized to perform this action",
			nil,
		)

	case apperrors.ErrDiscountRequestNotFound:
		response.Error(
			c,
			http.StatusNotFound,
			"discount request not found",
			nil,
		)

	case apperrors.ErrDiscountRequestNotPending:
		response.Error(
			c,
			http.StatusBadRequest,
			"discount request is not pending",
			nil,
		)

	default:
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to process discount request",
			err.Error(),
		)
	}
}
