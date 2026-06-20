package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/RitvikIP27/KubernesDeployment/database"
	"github.com/RitvikIP27/KubernesDeployment/models"
)

func GetDashboard(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var dash models.Dashboard
	database.DB.QueryRow("SELECT COUNT(*) FROM skills WHERE user_id = $1", userID).Scan(&dash.TotalSkills)
	database.DB.QueryRow("SELECT COALESCE(SUM(hours), 0) FROM learning_logs WHERE user_id = $1", userID).Scan(&dash.TotalHours)
	database.DB.QueryRow("SELECT COUNT(*) FROM learning_logs WHERE user_id = $1", userID).Scan(&dash.TotalLogs)

	err := database.DB.QueryRow(`
		SELECT s.name FROM skills s
		LEFT JOIN learning_logs l ON s.id = l.skill_id AND l.user_id = s.user_id
		WHERE s.user_id = $1
		GROUP BY s.id, s.name
		ORDER BY COALESCE(SUM(l.hours), 0) DESC
		LIMIT 1
	`, userID).Scan(&dash.TopSkill)
	if err != nil {
		dash.TopSkill = "N/A"
	}

	c.JSON(http.StatusOK, dash)
}

func HealthCheck(c *gin.Context) {
	err := database.DB.Ping()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
