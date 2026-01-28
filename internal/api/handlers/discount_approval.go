package handlers

import (
	"net/http"
	"strconv"

	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/response"

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
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid discount request id", nil)
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil && err.Error() != "EOF" {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	comment, _ := body["comment"].(string)

	err = services.ApproveDiscount(role, approverID, requestID, comment)
	if err != nil {
		handleApproveRejectDiscountError(c, err)
		return
	}

	response.Success(
		c,
		"discount request approved successfully",
		nil,
	)
}

func RejectDiscount(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid discount request id", nil)
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	comment, ok := body["comment"].(string)
	if !ok || comment == "" {
		response.Error(c, http.StatusBadRequest, "comment is required", nil)
		return
	}

	err = services.RejectDiscount(role, approverID, requestID, comment)
	if err != nil {
		handleApproveRejectDiscountError(c, err)
		return
	}

	response.Success(
		c,
		"discount request rejected successfully",
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
