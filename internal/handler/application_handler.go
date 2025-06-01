package handler

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
	"startfront-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateApplication handles the creation of a new application
func CreateApplication(c *gin.Context) {
	var app domain.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	err := usecase.CreateApplication(app)
	if err != nil {
		if err.Error() == "duplicate_found" {
			response.Error(c, "Duplicate code or route found")
		} else {
			response.Error(c, "Failed to create application")
		}
		return
	}

	// Send WebSocket message
	SendToClients("New application created")

	response.Success(c, "Application created successfully", nil)
}

// GetApplication fetches application details by code
func GetApplication(c *gin.Context) {
	id := c.Param("id")

	app, err := usecase.GetApplication(id)
	if err != nil {
		response.Error(c, "Application not found")
		return
	}

	response.Success(c, "", gin.H{"application": app})
}

// UpdateApplication updates application details
func UpdateApplication(c *gin.Context) {
	id := c.Param("id")
	var app domain.Application

	_, err := usecase.GetApplication(id)
	if err != nil {
		response.Error(c, "Application not found")
		return
	}

	if err := c.ShouldBindJSON(&app); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	err = usecase.UpdateApplication(id, app)
	if err != nil {
		response.Error(c, "Failed to update application")
		return
	}

	// Send WebSocket message
	SendToClients("Application updated")

	response.Success(c, "Application updated successfully", nil)
}

// DeleteApplication deletes an application
func DeleteApplication(c *gin.Context) {
	id := c.Param("id")

	_, err := usecase.GetApplication(id)
	if err != nil {
		response.Error(c, "Application not found")
		return
	}

	err = usecase.DeleteApplication(id)
	if err != nil {
		response.Error(c, "Failed to delete application")
		return
	}

	// Send WebSocket message
	SendToClients("Application deleted")

	response.Success(c, "Application deleted successfully", nil)
}

// ListApplications retrieves all applications
func ListApplications(c *gin.Context) {
	applications, err := usecase.ListApplications()
	if err != nil {
		response.Error(c, "Failed to fetch applications")
		return
	}

	response.Success(c, "", gin.H{
		"applications": applications,
	})
}
