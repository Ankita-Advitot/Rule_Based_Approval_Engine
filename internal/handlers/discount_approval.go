package handlers

import (
	"strconv"

	"rule-based-approval-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func GetPendingDiscounts(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	discounts, err := services.GetPendingDiscountRequests(role, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, discounts)
}

func ApproveDiscount(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request id"})
		return
	}

	err = services.ApproveDiscount(role, approverID, requestID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "discount approved"})
}

func RejectDiscount(c *gin.Context) {
	role := c.GetString("role")
	approverID := c.GetInt64("user_id")

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request id"})
		return
	}

	err = services.RejectDiscount(role, approverID, requestID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "discount rejected"})
}
