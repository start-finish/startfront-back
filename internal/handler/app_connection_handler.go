package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
)

func CreateAppConnection(c *gin.Context) {
	var conn domain.AppConnection
	if err := c.ShouldBindJSON(&conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid request", "error": err.Error()})
		return
	}

	if err := usecase.CreateAppConnection(conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Failed to create app connection", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "App connection created successfully"})
}

func GetAppConnectionsByAppID(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid app ID"})
		return
	}

	conns, err := usecase.ListAppConnectionsByAppID(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to fetch app connections", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": conns})
}
