package main

import (
	"log"
	"time"

	"rule-based-approval-engine/internal/config"
	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/jobs"
	"rule-based-approval-engine/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)

	router := gin.Default()

	//CORS MUST COME BEFORE ROUTES
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Register routes AFTER middleware
	routes.Register(router)

	//CRON SETUP
	loc, _ := time.LoadLocation("Asia/Kolkata")
	c := cron.New(cron.WithLocation(loc))
	c.AddFunc("0 0 * * *", jobs.RunAutoRejectJob)
	c.Start()

	log.Println("🚀 Server started on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
