package handlers

import (
	"net/http"
	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func RunAutoReject(c *gin.Context) {
	err1 := services.AutoRejectLeaveRequests()
	err2 := services.AutoRejectExpenseRequests()
	err3 := services.AutoRejectDiscountRequests()

	if err1 != nil || err2 != nil || err3 != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"auto reject failed",
			"one or more rejection processes encountered an error",
		)
		return
	}

	response.Success(
		c,
		"auto reject executed successfully",
		nil,
	)
}
