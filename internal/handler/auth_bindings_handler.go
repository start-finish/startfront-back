package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
)

func CreateAuthBinding(c *gin.Context) {
	var binding domain.AuthBinding
	if err := c.ShouldBindJSON(&binding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid request", "error": err.Error()})
		return
	}

	if err := usecase.CreateAuthBinding(binding); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to create auth binding", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Auth binding created successfully"})
}

func GetAuthBindingsByAppID(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("application_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid application ID"})
		return
	}

	bindings, err := usecase.ListAuthBindingsByAppID(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to get bindings", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": bindings})
}
