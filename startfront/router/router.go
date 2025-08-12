package router

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"startfront/handlers"
)

type Request struct {
	MsgID string          `json:"msgId"`
	Data  json.RawMessage `json:"data"`
}

type MsgHandler func(*gorm.DB, json.RawMessage, *gin.Context)

// Keep empty map on init; create-api will inject entries and add handlers import
var MsgHandlers = map[string]MsgHandler{
	"LOGIN": handlers.HandleLogin,
	"HOME": handlers.HandleHome,
}

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/api/doProcess", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "Invalid JSON",
			})
			return
		}
		handler, exists := MsgHandlers[req.MsgID]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "Unknown msgId",
			})
			return
		}
		handler(db, req.Data, c)
	})
}
