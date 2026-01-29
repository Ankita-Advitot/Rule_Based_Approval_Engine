package handlers

import (
	"net/http"
	"strconv"

	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func CreateRule(c *gin.Context) {
	role := c.GetString("role")

	var rule models.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := services.CreateRule(role, rule); err != nil {
		handleRuleError(c, err, "Failed to create rule")
		return
	}

	response.Created(c, "Rule created successfully", nil)
}

func GetRules(c *gin.Context) {
	role := c.GetString("role")

	rules, err := services.GetRules(role)
	if err != nil {
		handleRuleError(c, err, "Failed to fetch rules")
		return
	}

	response.Success(c, "Rules fetched successfully", rules)
}

func UpdateRule(c *gin.Context) {
	role := c.GetString("role")

	ruleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid rule ID", err.Error())
		return
	}

	var rule models.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := services.UpdateRule(role, ruleID, rule); err != nil {
		handleRuleError(c, err, "Failed to update rule")
		return
	}

	response.Success(c, "Rule updated successfully", nil)
}

func DeleteRule(c *gin.Context) {
	role := c.GetString("role")

	ruleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid rule ID", err.Error())
		return
	}

	if err := services.DeleteRule(role, ruleID); err != nil {
		handleRuleError(c, err, "Failed to delete rule")
		return
	}

	response.Success(c, "Rule deleted successfully", nil)
}

func handleRuleError(c *gin.Context, err error, message string) {
	status := http.StatusInternalServerError

	switch err {
	case apperrors.ErrUnauthorized:
		status = http.StatusForbidden
	case apperrors.ErrNoRuleFound, apperrors.ErrRuleNotFoundForDelete:
		status = http.StatusNotFound
	case apperrors.ErrRequestTypeRequired, apperrors.ErrActionRequired,
		apperrors.ErrGradeIDRequired, apperrors.ErrConditionRequired,
		apperrors.ErrInvalidConditionJSON:
		status = http.StatusBadRequest
	}

	response.Error(c, status, message, err.Error())
}
