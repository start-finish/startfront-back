package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleLogin(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	// Example: validate JSON body (optional)
	var payload map[string]interface{}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid request data: " + err.Error(),
			})
			return
		}
	}

	// TODO: add your DB logic with 'db' here

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "LOGIN API working",
	})
}
