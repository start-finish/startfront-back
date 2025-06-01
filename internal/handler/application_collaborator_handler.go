package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
)

func CreateApplicationCollaborator(c *gin.Context) {
	var collaborator domain.ApplicationCollaborator
	if err := c.ShouldBindJSON(&collaborator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid request"})
		return
	}

	err := usecase.CreateApplicationCollaborator(collaborator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Failed to create collaborator"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Collaborator added successfully"})
}

// GetApplicationCollaborators retrieves all collaborators for a given application
func GetApplicationCollaborators(c *gin.Context) {
	appIDStr := c.Param("application_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "Invalid application ID"})
		return
	}

	// Call usecase to fetch collaborators
	collaborators, err := usecase.ListApplicationCollaborators(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to fetch collaborators"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": collaborators})
}