package handlers

import (
	"net/http"
	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetMyLeaves(c *gin.Context) {
	userID := c.GetInt64("user_id")

	data, err := services.GetMyLeaveRequests(userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch leave requests",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"leave requests fetched successfully",
		data,
	)
}
func GetMyExpenses(c *gin.Context) {
	userID := c.GetInt64("user_id")

	data, err := services.GetMyExpenseRequests(userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch expense requests",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"expense requests fetched successfully",
		data,
	)
}
func GetMyDiscounts(c *gin.Context) {
	userID := c.GetInt64("user_id")

	data, err := services.GetMyDiscountRequests(userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch discount requests",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"discount requests fetched successfully",
		data,
	)
}
