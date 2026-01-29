package main

import (
	"log"
	"time"

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

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Register(router)

	//CRON SETUP
	loc, _ := time.LoadLocation("Asia/Kolkata")
	c := cron.New(cron.WithLocation(loc))
	c.AddFunc("0 0 * * *", jobs.RunAutoRejectJob)
	// c.AddFunc("every @1m", jobs.RunAutoRejectJob)

	c.Start()

	log.Println("🚀 Server started on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
