package routes

import (
	"rule-based-approval-engine/internal/api/handlers"
	"rule-based-approval-engine/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", handlers.Logout)

	}

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetInt64("user_id"),
				"role":    c.GetString("role"),
			})
		})

		// Balances
		protected.GET("/balances", handlers.GetMyBalances)

		leaves := protected.Group("/leaves")
		{
			leaves.POST("/request", handlers.ApplyLeave)
			leaves.POST("/:id/cancel", handlers.CancelLeave)
			leaves.GET("/my", handlers.GetMyLeaves)

			leaves.GET("/pending", handlers.GetPendingLeaves)
			leaves.POST("/:id/approve", handlers.ApproveLeave)
			leaves.POST("/:id/reject", handlers.RejectLeave)
		}

		expenses := protected.Group("/expenses")
		{
			expenses.POST("/request", handlers.ApplyExpense)
			expenses.POST("/:id/cancel", handlers.CancelExpense)
			expenses.GET("/my", handlers.GetMyExpenses)

			// Manager/Admin (if you add later)
			expenses.GET("/pending", handlers.GetPendingExpenses)
			expenses.POST("/:id/approve", handlers.ApproveExpense)
			expenses.POST("/:id/reject", handlers.RejectExpense)
		}

		discounts := protected.Group("/discounts")
		{
			discounts.POST("/request", handlers.ApplyDiscount)
			discounts.POST("/:id/cancel", handlers.CancelDiscount)
			discounts.GET("/my", handlers.GetMyDiscounts)

			// Manager/Admin (if you add later)

			discounts.GET("/pending", handlers.GetPendingDiscounts)
			discounts.POST("/:id/approve", handlers.ApproveDiscount)
			discounts.POST("/:id/reject", handlers.RejectDiscount)
		}

		system := protected.Group("/system")
		{

			// Manual trigger for testing auto-reject
			system.POST("/run-auto-reject", handlers.RunAutoReject)
		}

		admin := protected.Group("/admin")
		{
			admin.POST("/holidays", handlers.AddHoliday)
			admin.GET("/holidays", handlers.GetHolidays)
			admin.DELETE("/holidays/:id", handlers.DeleteHoliday)
			admin.POST("/rules", handlers.CreateRule)
			admin.GET("/rules", handlers.GetRules)
			admin.PUT("/rules/:id", handlers.UpdateRule)
			admin.DELETE("/rules/:id", handlers.DeleteRule)

			admin.GET("/reports/request-status-distribution", handlers.GetRequestStatusDistribution)
			admin.GET("/reports/requests-by-type", handlers.GetRequestsByType)
		}

	}
}
