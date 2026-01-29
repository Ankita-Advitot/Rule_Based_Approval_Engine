package handlers

import (
	"net/http"

	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid input",
			err.Error(),
		)
		return
	}

	err := services.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		status := http.StatusBadRequest
		if err == apperrors.ErrEmailAlreadyRegistered {
			status = http.StatusConflict
		}
		response.Error(
			c,
			status,
			"registration failed",
			err.Error(),
		)
		return
	}

	response.Created(
		c,
		"user registered successfully",
		nil,
	)
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid request", err.Error())
		return
	}

	token, role, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err == apperrors.ErrEmailRequired || err == apperrors.ErrPasswordRequired {
			status = http.StatusBadRequest
		}
		response.Error(c, status, "invalid credentials", err.Error())
		return
	}

	c.SetCookie(
		"access_token",
		token,
		3600*24, // 1 day
		"/",
		"",
		false,
		true,
	)

	response.Success(
		c,
		"login successful",
		gin.H{
			"token": token,
			"role":  role,
		},
	)
}

func Logout(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	response.Success(
		c,
		"logged out successfully",
		nil,
	)
}
