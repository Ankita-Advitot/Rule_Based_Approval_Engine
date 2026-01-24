package handlers

import (
	"context"
	"log"
	"net/http"

	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/response"

	"github.com/gin-gonic/gin"
)

func GetMyBalances(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var leaveTotal, leaveRemaining int
	var expenseTotal, expenseRemaining float64
	var discountTotal, discountRemaining float64

	err := database.DB.QueryRow(
		context.Background(),
		`SELECT total_allocated, remaining_count FROM leaves WHERE user_id=$1`,
		userID,
	).Scan(&leaveTotal, &leaveRemaining)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch leave balance",
			err.Error(),
		)
		log.Printf("Error fetching leave balance: %v", err)
		return
	}
	err = database.DB.QueryRow(
		context.Background(),
		`SELECT total_amount, remaining_amount FROM expense WHERE user_id=$1`,
		userID,
	).Scan(&expenseTotal, &expenseRemaining)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch expense balance",
			err.Error(),
		)
		return
	}

	err = database.DB.QueryRow(
		context.Background(),
		`SELECT total_discount, remaining_discount FROM discount WHERE user_id=$1`,
		userID,
	).Scan(&discountTotal, &discountRemaining)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch discount balance",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"balances fetched successfully",
		gin.H{
			"leave": gin.H{
				"total":     leaveTotal,
				"remaining": leaveRemaining,
			},
			"expense": gin.H{
				"total":     expenseTotal,
				"remaining": expenseRemaining,
			},
			"discount": gin.H{
				"total":     discountTotal,
				"remaining": discountRemaining,
			},
		},
	)
}
