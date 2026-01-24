package handlers

import (
	"net/http"
	"strconv"

	"rule-based-approval-engine/internal/response"
	"rule-based-approval-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func GetPendingLeaves(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	leaves, err := services.GetPendingLeaveRequests(role, userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch pending leave requests",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"pending leave requests fetched successfully",
		leaves,
	)
}
func ApproveLeave(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid leave request id",
			nil,
		)
		return
	}

	err = services.ApproveLeave(role, approverID, requestID)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"unable to approve leave request",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"leave approved successfully",
		nil,
	)
}

func RejectLeave(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid leave request id",
			nil,
		)
		return
	}

	err = services.RejectLeave(role, approverID, requestID)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"unable to reject leave request",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"leave rejected successfully",
		nil,
	)
}
