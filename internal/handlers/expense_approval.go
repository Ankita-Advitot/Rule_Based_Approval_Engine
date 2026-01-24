package handlers

import (
	"strconv"

	"rule-based-approval-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func GetPendingExpenses(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	expenses, err := services.GetPendingExpenseRequests(role, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, expenses)
}

func ApproveExpense(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request id"})
		return
	}

	err = services.ApproveExpense(role, approverID, requestID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "expense approved"})
}

func RejectExpense(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request id"})
		return
	}

	err = services.RejectExpense(role, approverID, requestID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "expense rejected"})
}
