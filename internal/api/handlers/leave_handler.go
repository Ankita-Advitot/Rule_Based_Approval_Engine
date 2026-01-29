package handlers

import (
	"net/http"
	"strconv"
	"time"

	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/response"
	"rule-based-approval-engine/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type LeaveApplyRequest struct {
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	LeaveType string `json:"leave_type"`
	Reason    string `json:"reason"`
}

func ApplyLeave(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "unauthorized user", nil)
		return
	}

	var req LeaveApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	from, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid from_date format (YYYY-MM-DD)", nil)
		return
	}

	to, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid to_date format (YYYY-MM-DD)", nil)
		return
	}

	days := utils.CalculateLeaveDays(from, to)
	if days <= 0 {
		response.Error(c, http.StatusBadRequest, "leave duration must be at least one day", nil)
		return
	}

	message, status, err := services.ApplyLeave(
		userID, from, to, days, req.LeaveType, req.Reason,
	)

	if err != nil {
		handleApplyLeaveError(c, err)
		return
	}

	response.Success(c, message, gin.H{
		"status": status,
	})

}

func handleApplyLeaveError(c *gin.Context, err error) {
	switch err {

	case apperrors.ErrLeaveBalanceExceeded:
		response.Error(c, http.StatusBadRequest, "insufficient leave balance", nil)

	case apperrors.ErrInvalidLeaveDays:
		response.Error(c, http.StatusBadRequest, "invalid leave duration", nil)

	case apperrors.ErrRuleNotFound:
		response.Error(c, http.StatusInternalServerError, "leave approval rules not configured", nil)

	case apperrors.ErrUserNotFound:
		response.Error(c, http.StatusNotFound, "user not found", nil)

	case apperrors.ErrLeaveBalanceMissing:
		response.Error(c, http.StatusNotFound, "leave balance not found", nil)

	case apperrors.ErrLeaveOverlap:
		response.Error(c, http.StatusBadRequest, "overlapping leave request exists", nil)

	default:
		response.Error(c, http.StatusInternalServerError, "failed to apply leave", err.Error())
	}
}

func CancelLeave(c *gin.Context) {
	userID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid leave request id", nil)
		return
	}

	err = services.CancelLeave(userID, requestID)
	if err != nil {
		handleCancelLeaveError(c, err)
		return
	}

	response.Success(c, "leave request cancelled successfully", nil)
}

func handleCancelLeaveError(c *gin.Context, err error) {
	status := http.StatusInternalServerError

	switch err {
	case apperrors.ErrLeaveRequestNotFound:
		status = http.StatusNotFound
	case apperrors.ErrRequestCannotCancel:
		status = http.StatusBadRequest
	}

	response.Error(c, status, "unable to cancel leave", err.Error())
}
