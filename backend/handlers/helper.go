package handlers

import "github.com/gin-gonic/gin"

func getCurrentUserID(c *gin.Context) (int, bool) {
	id, ok := c.Get("currentUserID")
	if !ok {
		return 0, false
	}

	userID, ok := id.(int)
	return userID, ok
}
