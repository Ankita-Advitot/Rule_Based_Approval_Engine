package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ankita-advitot/rule_based_approval_engine/app/auth"
	"github.com/ankita-advitot/rule_based_approval_engine/app/auto_reject"
	"github.com/ankita-advitot/rule_based_approval_engine/app/domain_service"
	"github.com/ankita-advitot/rule_based_approval_engine/app/expense_service"
	"github.com/ankita-advitot/rule_based_approval_engine/app/holidays"
	"github.com/ankita-advitot/rule_based_approval_engine/app/leave_service"
	"github.com/ankita-advitot/rule_based_approval_engine/app/my_requests"
	"github.com/ankita-advitot/rule_based_approval_engine/app/reports"
	"github.com/ankita-advitot/rule_based_approval_engine/app/rules"
	"github.com/ankita-advitot/rule_based_approval_engine/config"
	jobs "github.com/ankita-advitot/rule_based_approval_engine/cron-jobs"
	"github.com/ankita-advitot/rule_based_approval_engine/database"
	"github.com/ankita-advitot/rule_based_approval_engine/repositories"
	"github.com/ankita-advitot/rule_based_approval_engine/repositories/migrations"
	"github.com/ankita-advitot/rule_based_approval_engine/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrations.HandleMigrateCommand()
		return
	}

	ctx := context.Background()

	// 1. Core Repositories
	userRepo := repositories.NewUserRepository(ctx, database.DB)
	balanceRepo := repositories.NewBalanceRepository(ctx, database.DB)
	ruleRepo := repositories.NewRuleRepository(ctx, database.DB)
	leaveRepo := repositories.NewLeaveRequestRepository(ctx, database.DB)
	expenseRepo := repositories.NewExpenseRequestRepository(ctx, database.DB)
	discountRepo := repositories.NewDiscountRequestRepository(ctx, database.DB)
	holidayRepo := repositories.NewHolidayRepository(ctx, database.DB)
	reportRepo := repositories.NewReportRepository(ctx, database.DB)
	myRequestsRepo := repositories.NewAggregatedRepository(ctx, database.DB)

	// 2. Services
	authService := auth.NewAuthService(ctx, userRepo, balanceRepo, database.DB)
	ruleService := rules.NewRuleService(ctx, ruleRepo)
	leaveService := leave_service.NewLeaveService(
		ctx, leaveRepo, balanceRepo, ruleService, userRepo, database.DB,
	)
	leaveApprovalService := leave_service.NewLeaveApprovalService(
		ctx, leaveRepo, balanceRepo, userRepo, database.DB,
	)
	expenseService := expense_service.NewExpenseService(
		ctx, expenseRepo, balanceRepo, ruleService, userRepo, database.DB,
	)
	expenseApprovalService := expense_service.NewExpenseApprovalService(
		ctx, expenseRepo, balanceRepo, userRepo, database.DB,
	)
	holidayService := holidays.NewHolidayService(ctx, holidayRepo)
	reportService := reports.NewReportService(ctx, reportRepo)
	balanceService := domain_service.NewBalanceService(ctx, balanceRepo, database.DB)
	discountService := domain_service.NewDiscountService(ctx, discountRepo, balanceRepo, ruleService, userRepo, database.DB)
	discountApprovalService := domain_service.NewDiscountApprovalService(ctx, discountRepo, balanceRepo, userRepo, database.DB)
	autoRejectService := auto_reject.NewAutoRejectService(
		ctx, leaveRepo, expenseRepo, discountRepo, holidayRepo, database.DB,
	)
	myRequestsService := my_requests.NewMyRequestsService(ctx, myRequestsRepo)

	// 3. Router & CORS
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 4. Route Registration
	routes.Register(
		ctx,
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
	)

	// 5. Cron Jobs
	loc, _ := time.LoadLocation("Asia/Kolkata")
	c := cron.New(cron.WithLocation(loc))
	c.AddFunc("0 0 * * *", func() {
		jobs.RunAutoRejectJob(ctx, autoRejectService)
	})
	c.Start()

	log.Println(" Server started on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
