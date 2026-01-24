package handlers

import (
	"net/http"
	"strconv"

	"rule-based-approval-engine/internal/response"
	"rule-based-approval-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func GetPendingExpenses(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	expenses, err := services.GetPendingExpenseRequests(role, userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch pending expense requests",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"pending expense requests fetched successfully",
		expenses,
	)
}

func ApproveExpense(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

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

	err = services.ApproveExpense(role, approverID, requestID)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"unable to approve expense request",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"expense approved successfully",
		nil,
	)
}

func RejectExpense(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

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

	err = services.RejectExpense(role, approverID, requestID)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"unable to reject expense request",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"expense rejected successfully",
		nil,
	)
}
