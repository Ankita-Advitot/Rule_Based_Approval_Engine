package handlers

import (
	"net/http"

	"rule-based-approval-engine/internal/response"
	"rule-based-approval-engine/internal/services"

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
		response.Error(
			c,
			http.StatusBadRequest,
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid input",
			err.Error(),
		)
		return
	}

	token, role, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"login failed",
			err.Error(),
		)
		return
	}

	// Set JWT as HttpOnly cookie
	c.SetCookie(
		"access_token",
		token,
		3600*24,
		"/",
		"",
		false,
		true,
	)

	response.Success(
		c,
		"login successful",
		gin.H{
			"role": role,
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
