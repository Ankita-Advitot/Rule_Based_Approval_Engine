package main

import (
	"log"
	"time"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/app/services"
	"rule-based-approval-engine/internal/config"
	jobs "rule-based-approval-engine/internal/cron-jobs"
	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)

	userRepo := repositories.NewUserRepository(database.DB)
	balanceRepo := repositories.NewBalanceRepository(database.DB)
	authService := services.NewAuthService(userRepo, balanceRepo, database.DB)

	ruleRepo := repositories.NewRuleRepository(database.DB)
	ruleService := services.NewRuleService(ruleRepo)

	leaveRepo := repositories.NewLeaveRequestRepository(database.DB)
	leaveService := services.NewLeaveService(
		leaveRepo, balanceRepo, ruleService, userRepo, database.DB,
	)
	leaveApprovalService := services.NewLeaveApprovalService(
		leaveRepo, balanceRepo, userRepo, database.DB,
	)

	expenseRepo := repositories.NewExpenseRequestRepository(database.DB)
	expenseService := services.NewExpenseService(
		expenseRepo, balanceRepo, ruleService, userRepo, database.DB,
	)
	expenseApprovalService := services.NewExpenseApprovalService(
		expenseRepo, balanceRepo, userRepo, database.DB,
	)

	myRequestsRepo := repositories.NewMyRequestsRepository(database.DB)
	myRequestsService := services.NewMyRequestsService(myRequestsRepo)

	holidayRepo := repositories.NewHolidayRepository(database.DB)
	holidayService := services.NewHolidayService(holidayRepo)

	reportRepo := repositories.NewReportRepository(database.DB)
	reportService := services.NewReportService(reportRepo)

	balanceService := services.NewBalanceService(balanceRepo, database.DB)

	discountRepo := repositories.NewDiscountRequestRepository(database.DB)
	discountService := services.NewDiscountService(discountRepo, balanceRepo, ruleService, userRepo, database.DB)
	discountApprovalService := services.NewDiscountApprovalService(discountRepo, balanceRepo, userRepo, database.DB)

	autoRejectService := services.NewAutoRejectService(
		leaveRepo, expenseRepo, discountRepo, holidayRepo, database.DB,
	)

	// Aggregated My Requests
	aggregatedRepo := repositories.NewAggregatedRepository(database.DB)
	myRequestsAllService := services.NewMyRequestsServices(aggregatedRepo)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Register(
		router,
		authService,
		leaveService,
		expenseService,
		leaveApprovalService,
		expenseApprovalService,
		ruleService,
		myRequestsService,
		holidayService,
		reportService,
		balanceService,
		autoRejectService,
		discountService,
		discountApprovalService,
		myRequestsAllService,
	)

	loc, _ := time.LoadLocation("Asia/Kolkata")
	c := cron.New(cron.WithLocation(loc))
	c.AddFunc("0 0 * * *", func() {
		jobs.RunAutoRejectJob(autoRejectService)
	})
	c.Start()

	log.Println(" Server started on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
