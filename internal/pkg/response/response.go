package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(
		200,
		gin.H{
			"success": true,
			"message": message,
			"data":    data,
			"error":   nil,
		},
	)
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(
		201,
		gin.H{
			"success": true,
			"message": message,
			"data":    data,
			"error":   nil,
		},
	)
}

func Error(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(
		status,
		gin.H{
			"success": false,
			"message": message,
			"data":    nil,
			"error":   err,
		},
	)
}
