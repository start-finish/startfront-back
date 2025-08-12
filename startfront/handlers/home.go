package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"startfront/models"
)

func HandleHome(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	// Step 1: Parse incoming JSON into a struct for safety
	var req struct {
		ID uint `json:"id"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "invalid request data: " + err.Error(),
			})
			return
		}
	}

	// Validate ID
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   "1",
			"status": "error",
			"error":  "'id' must be greater than 0",
		})
		return
	}

	// Step 2: Query the DB for the record
	var home models.Home
	if err := db.First(&home, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":   "1",
				"status": "error",
				"error":  "record not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   "1",
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Step 3: Return the record
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "HOME API working",
		"data":    home,
	})
}
