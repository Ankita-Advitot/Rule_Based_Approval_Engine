package middleware

import (
	"strings"

	"rule-based-approval-engine/internal/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// if tokenString == "" {
		cookieToken, err := c.Cookie("access_token")
		if err == nil {
			tokenString = cookieToken
		}
		// }

		//  If still empty → unauthorized
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
