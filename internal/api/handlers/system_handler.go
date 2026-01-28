package handlers

import (
	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/response"

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
