package routes

import (
	"rule-based-approval-engine/internal/api/handlers"
	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Register(
	r *gin.Engine,
	authService *services.AuthService,
	leaveService *services.LeaveService,
	expenseService *services.ExpenseService,
	leaveApprovalService *services.LeaveApprovalService,
	expenseApprovalService *services.ExpenseApprovalService,
	ruleService *services.RuleService,
	myRequestsService *services.MyRequestsService,
	holidayService *services.HolidayService,
	reportService *services.ReportService,
	balanceService *services.BalanceService,
	autoRejectService *services.AutoRejectService,
	discountService *services.DiscountService,
	discountApprovalService services.DiscountApprovalServiceInterface,
	myRequestsAllService services.MyRequestsServices,
) {
	// Initialize handlers with services
	authHandler := handlers.NewAuthHandler(authService)
	leaveHandler := handlers.NewLeaveHandler(leaveService)
	leaveApprovalHandler := handlers.NewLeaveApprovalHandler(leaveApprovalService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	expenseApprovalHandler := handlers.NewExpenseApprovalHandler(expenseApprovalService)
	ruleHandler := handlers.NewRuleHandler(ruleService)
	myRequestsHandler := handlers.NewMyRequestsHandler(myRequestsService)
	holidayHandler := handlers.NewHolidayHandler(holidayService)
	reportHandler := handlers.NewReportHandler(reportService)
	balanceHandler := handlers.NewBalanceHandler(balanceService)
	systemHandler := handlers.NewSystemHandler(autoRejectService)
	discountHandler := handlers.NewDiscountHandler(discountService)
	discountApprovalHandler := handlers.NewDiscountApprovalHandler(discountApprovalService)
	myRequestsAllHandler := handlers.NewMyRequestsHandlers(myRequestsAllService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)

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

		// Aggregated My Requests
		protected.GET("/my_requests", myRequestsAllHandler.GetMyAllRequests)

		// Balances
		protected.GET("/balances", balanceHandler.GetMyBalances)

		leaves := protected.Group("/leaves")
		{
			leaves.POST("/request", leaveHandler.ApplyLeave)
			leaves.POST("/:id/cancel", leaveHandler.CancelLeave)
			leaves.GET("/my", myRequestsHandler.GetMyLeaves)

			approval := leaves.Group(
				"",
				middleware.RequireRole("MANAGER", "ADMIN"),
			)
			{
				approval.GET("/pending", leaveApprovalHandler.GetPendingLeaves)
				approval.POST("/:id/approve", leaveApprovalHandler.ApproveLeave)
				approval.POST("/:id/reject", leaveApprovalHandler.RejectLeave)
			}
		}

		expenses := protected.Group("/expenses")
		{
			expenses.POST("/request", expenseHandler.ApplyExpense)
			expenses.POST("/:id/cancel", expenseHandler.CancelExpense)
			expenses.GET("/my", myRequestsHandler.GetMyExpenses)

			approval := expenses.Group(
				"",
				middleware.RequireRole("MANAGER", "ADMIN"),
			)
			{
				approval.GET("/pending", expenseApprovalHandler.GetPendingExpenses)
				approval.POST("/:id/approve", expenseApprovalHandler.ApproveExpense)
				approval.POST("/:id/reject", expenseApprovalHandler.RejectExpense)
			}
		}
		discounts := protected.Group("/discounts")
		{
			discounts.POST("/request", discountHandler.ApplyDiscount)
			discounts.POST("/:id/cancel", discountHandler.CancelDiscount)
			discounts.GET("/my", myRequestsHandler.GetMyDiscounts)

			approval := discounts.Group(
				"",
				middleware.RequireRole("MANAGER", "ADMIN"),
			)
			{
				approval.GET("/pending", discountApprovalHandler.GetPendingDiscounts)
				approval.POST("/:id/approve", discountApprovalHandler.ApproveDiscount)
				approval.POST("/:id/reject", discountApprovalHandler.RejectDiscount)
			}
		}

		system := protected.Group("/system", middleware.RequireRole("ADMIN"))
		{

			// Manual trigger for testing auto-reject
			system.POST("/run-auto-reject", systemHandler.RunAutoReject)
		}

		admin := protected.Group("/admin", middleware.RequireRole("ADMIN"))
		{
			admin.POST("/holidays", holidayHandler.AddHoliday)
			admin.GET("/holidays", holidayHandler.GetHolidays)
			admin.DELETE("/holidays/:id", holidayHandler.DeleteHoliday)
			admin.POST("/rules", ruleHandler.CreateRule)
			admin.GET("/rules", ruleHandler.GetRules)
			admin.PUT("/rules/:id", ruleHandler.UpdateRule)
			admin.DELETE("/rules/:id", ruleHandler.DeleteRule)

			admin.GET("/reports/request-status-distribution", reportHandler.GetRequestStatusDistribution)
			admin.GET("/reports/requests-by-type", reportHandler.GetRequestsByType)
		}

	}
}
