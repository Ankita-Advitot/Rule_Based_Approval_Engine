package handlers

import (
	"rule-based-approval-engine/internal/response"
	"rule-based-approval-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func RunAutoReject(c *gin.Context) {
	services.AutoRejectLeaveRequests()
	services.AutoRejectExpenseRequests()
	services.AutoRejectDiscountRequests()

	response.Success(
		c,
		"auto reject executed successfully",
		nil,
	)
}
