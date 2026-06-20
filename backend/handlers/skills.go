package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/RitvikIP27/KubernesDeployment/database"
	"github.com/RitvikIP27/KubernesDeployment/models"
)

func GetSkills(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT s.id, s.name, s.category, s.target_hours,
		       COALESCE(SUM(l.hours), 0) as total_hours, s.created_at
		FROM skills s
		LEFT JOIN learning_logs l ON s.id = l.skill_id AND l.user_id = s.user_id
		WHERE s.user_id = $1
		GROUP BY s.id, s.name, s.category, s.target_hours, s.created_at
		ORDER BY s.created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	skills := []models.Skill{}
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Category, &s.TargetHours, &s.TotalHours, &s.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		skills = append(skills, s)
	}

	c.JSON(http.StatusOK, skills)
}

func CreateSkill(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	err := database.DB.QueryRow(
		"INSERT INTO skills (user_id, name, category, target_hours) VALUES ($1, $2, $3, $4) RETURNING id",
		userID, req.Name, req.Category, req.TargetHours,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Skill created"})
}

func GetSkill(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	var skill models.Skill
	err := database.DB.QueryRow(`
		SELECT s.id, s.name, s.category, s.target_hours,
		       COALESCE(SUM(l.hours), 0) as total_hours, s.created_at
		FROM skills s
		LEFT JOIN learning_logs l ON s.id = l.skill_id AND l.user_id = s.user_id
		WHERE s.id = $1 AND s.user_id = $2
		GROUP BY s.id, s.name, s.category, s.target_hours, s.created_at
	`, id, userID).Scan(&skill.ID, &skill.Name, &skill.Category, &skill.TargetHours, &skill.TotalHours, &skill.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	rows, err := database.DB.Query(
		"SELECT id, skill_id, hours, notes, log_date, created_at FROM learning_logs WHERE skill_id = $1 AND user_id = $2 ORDER BY log_date DESC",
		id, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	logs := []models.LearningLog{}
	for rows.Next() {
		var l models.LearningLog
		if err := rows.Scan(&l.ID, &l.SkillID, &l.Hours, &l.Notes, &l.LogDate, &l.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs = append(logs, l)
	}

	c.JSON(http.StatusOK, gin.H{"skill": skill, "logs": logs})
}

func DeleteSkill(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM skills WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted"})
}
