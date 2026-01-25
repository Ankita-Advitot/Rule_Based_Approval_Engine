package handlers

import (
	"net/http"
	"strconv"

	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/response"
	"rule-based-approval-engine/internal/services"

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
		response.Error(c, http.StatusForbidden, "Failed to create rule", err.Error())
		return
	}

	response.Created(c, "Rule created successfully", nil)
}

func GetRules(c *gin.Context) {
	role := c.GetString("role")

	rules, err := services.GetRules(role)
	if err != nil {
		response.Error(c, http.StatusForbidden, "Failed to fetch rules", err.Error())
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
		response.Error(c, http.StatusForbidden, "Failed to update rule", err.Error())
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
		response.Error(c, http.StatusForbidden, "Failed to delete rule", err.Error())
		return
	}

	response.Success(c, "Rule deleted successfully", nil)
}
