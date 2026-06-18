package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trainwithshubham/skillpulse/database"
)

type UpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func GetSettings(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rows, err := database.DB.Query("SELECT key, value FROM settings WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, val string
		if err := rows.Scan(&key, &val); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		settings[key] = val
	}

	c.JSON(http.StatusOK, settings)
}

func UpdateSetting(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(
		`INSERT INTO settings (user_id, key, value) VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, key) DO UPDATE SET value = EXCLUDED.value`,
		userID, req.Key, req.Value,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully"})
}
