package handlers

import (
	"net/http"
	"strconv"
	"time"

	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type HolidayRequest struct {
	Date        string `json:"date"`
	Description string `json:"description"`
}

func AddHoliday(c *gin.Context) {
	role := c.GetString("role")
	adminID := c.GetInt64("user_id")

	var req HolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid input",
			err.Error(),
		)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid date format (YYYY-MM-DD)",
			nil,
		)
		return
	}

	err = services.AddHoliday(role, adminID, date, req.Description)
	if err != nil {
		response.Error(
			c,
			http.StatusForbidden,
			"unable to add holiday",
			err.Error(),
		)
		return
	}

	response.Created(
		c,
		"holiday added successfully",
		nil,
	)
}

func GetHolidays(c *gin.Context) {
	role := c.GetString("role")

	holidays, err := services.GetHolidays(role)
	if err != nil {
		response.Error(
			c,
			http.StatusForbidden,
			"failed to fetch holidays",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"holidays fetched successfully",
		holidays,
	)
}

func DeleteHoliday(c *gin.Context) {
	role := c.GetString("role")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid holiday id",
			nil,
		)
		return
	}

	err = services.DeleteHoliday(role, id)
	if err != nil {
		response.Error(
			c,
			http.StatusForbidden,
			"unable to delete holiday",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		"holiday removed successfully",
		nil,
	)
}
