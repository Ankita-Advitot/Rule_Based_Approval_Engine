package handlers

import (
	"net/http"

	"rule-based-approval-engine/internal/response"
	"rule-based-approval-engine/internal/services"

	"github.com/gin-gonic/gin"
)

func GetRequestStatusDistribution(c *gin.Context) {
	role := c.GetString("role")
	if role != "ADMIN" {
		response.Error(
			c,
			http.StatusForbidden,
			"Unauthorized access",
			"admin role required",
		)
		return
	}

	data, err := services.GetRequestStatusDistribution()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to fetch request status distribution",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"Request status distribution fetched successfully",
		data,
	)
}

func GetRequestsByType(c *gin.Context) {
	role := c.GetString("role")
	if role != "ADMIN" {
		response.Error(
			c,
			http.StatusForbidden,
			"Unauthorized access",
			"admin role required",
		)
		return
	}

	data, err := services.GetRequestsByTypeReport()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to fetch requests by type report",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"Requests by type report fetched successfully",
		data,
	)
}
