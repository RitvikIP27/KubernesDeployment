package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/trainwithshubham/skillpulse/database"
	"github.com/trainwithshubham/skillpulse/handlers"
	"github.com/trainwithshubham/skillpulse/middleware"
)

func main() {
	database.Connect()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing JWT_SECRET environment variable")
	}

	router := gin.Default()

	// Auth routes
	router.POST("/auth/register", handlers.Register)
	router.POST("/auth/login", handlers.Login)
	router.GET("/auth/me", middleware.AuthMiddleware(jwtSecret), handlers.GetMe)

	// API routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/skills", handlers.GetSkills)
		api.POST("/skills", handlers.CreateSkill)
		api.GET("/skills/:id", handlers.GetSkill)
		api.DELETE("/skills/:id", handlers.DeleteSkill)
		api.POST("/skills/:id/log", handlers.CreateLog)
		api.GET("/dashboard", handlers.GetDashboard)
		api.GET("/settings", handlers.GetSettings)
		api.POST("/settings", handlers.UpdateSetting)
	}

	// Health check
	router.GET("/health", handlers.HealthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SkillPulse API running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
