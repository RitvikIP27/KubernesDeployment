package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/RitvikIP27/KubernesDeployment/database"
	"github.com/RitvikIP27/KubernesDeployment/models"
)

func CreateLog(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skillID := c.Param("id")

	var exists bool
	err := database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM skills WHERE id = $1 AND user_id = $2)",
		skillID, userID,
	).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	var req models.CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	err = database.DB.QueryRow(
		"INSERT INTO learning_logs (skill_id, user_id, hours, notes, log_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		skillID, userID, req.Hours, req.Notes, req.LogDate,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Learning session logged"})
}
