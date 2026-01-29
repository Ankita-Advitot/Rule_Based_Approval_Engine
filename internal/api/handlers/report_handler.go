package handlers

import (
	"net/http"

	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetRequestStatusDistribution(c *gin.Context) {
	role := c.GetString("role")
	if role != "ADMIN" {
		handleReportError(c, apperrors.ErrAdminOnly, "Unauthorized access")
		return
	}

	data, err := services.GetRequestStatusDistribution()
	if err != nil {
		handleReportError(c, err, "Failed to fetch request status distribution")
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
		handleReportError(c, apperrors.ErrAdminOnly, "Unauthorized access")
		return
	}

	data, err := services.GetRequestsByTypeReport()
	if err != nil {
		handleReportError(c, err, "Failed to fetch requests by type report")
		return
	}

	response.Success(
		c,
		"Requests by type report fetched successfully",
		data,
	)
}

func handleReportError(c *gin.Context, err error, message string) {
	status := http.StatusInternalServerError

	if err == apperrors.ErrAdminOnly || err == apperrors.ErrUnauthorized {
		status = http.StatusForbidden
	}

	response.Error(c, status, message, err.Error())
}
