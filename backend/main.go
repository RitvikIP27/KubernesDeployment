package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/RitvikIP27/KubernesDeployment/database"
	"github.com/RitvikIP27/KubernesDeployment/handlers"
	"github.com/RitvikIP27/KubernesDeployment/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	_ = godotenv.Load("../.env")
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
	// router.GET("/auth/google", handlers.GoogleLogin)
	// router.GET("/auth/google/callback", handlers.GoogleCallback)

	authAPI := router.Group("/api/auth")
	{
		authAPI.POST("/register", handlers.Register)
		authAPI.POST("/login", handlers.Login)
		authAPI.GET("/me", middleware.AuthMiddleware(jwtSecret), handlers.GetMe)
		authAPI.GET("/google", handlers.GoogleLogin)
		authAPI.GET("/google/callback", handlers.GoogleCallback)
	}

	// API routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Skills & Learning
		api.GET("/skills", handlers.GetSkills)
		api.POST("/skills", handlers.CreateSkill)
		api.GET("/skills/:id", handlers.GetSkill)
		api.DELETE("/skills/:id", handlers.DeleteSkill)
		api.POST("/skills/:id/log", handlers.CreateLog)

		// Dashboard & Analytics (Phase 2)
		api.GET("/dashboard", handlers.GetDashboard)
		api.GET("/analytics", handlers.GetAnalytics)

		// Career Readiness (Phase 3)
		api.GET("/career-readiness", handlers.GetCareerReadiness)
		api.GET("/career-readiness/:track", handlers.GetCareerReadinessByTrack)

		// Job Preferences & Matching (Phase 4)
		api.GET("/job-preferences", handlers.GetJobPreferences)
		api.POST("/job-preferences", handlers.AddJobPreference)
		api.DELETE("/job-preferences/:role", handlers.RemoveJobPreference)
		api.GET("/job-matches", handlers.GetJobMatches)
		api.GET("/job-matches/:role", handlers.GetJobMatch)

		// Projects (Phase 5)
		api.GET("/projects", handlers.GetProjects)
		api.POST("/projects", handlers.CreateProject)
		api.GET("/projects/:id", handlers.GetProject)
		api.PUT("/projects/:id", handlers.UpdateProject)
		api.DELETE("/projects/:id", handlers.DeleteProject)

		// Certificates (Phase 5)
		api.GET("/certificates", handlers.GetCertificates)
		api.POST("/certificates", handlers.CreateCertificate)
		api.DELETE("/certificates/:id", handlers.DeleteCertificate)

		// User Profile & Professional Profile (Phase 5)
		api.GET("/profile", handlers.GetUserProfile)
		api.PUT("/profile", handlers.UpdateUserProfile)
		api.GET("/profile/professional", handlers.GenerateProfessionalProfile)

		// Settings
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
