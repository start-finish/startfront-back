package handler

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
	"startfront-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateScreen handles the creation of a new screen
func CreateScreen(c *gin.Context) {
	var screen domain.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	err := usecase.CreateScreen(screen)
	if err != nil {
		if err.Error() == "duplicate_found" {
			response.Error(c, "Duplicate code or route found")
		} else {
			response.Error(c, "Failed to create screen")
		}
		return
	}

	// Send WebSocket message
	SendToClients("New screen created")

	response.Success(c, "Screen created successfully", nil)
}

// GetScreensById fetches screens by application code
func GetScreensById(c *gin.Context) {
	id := c.Param("id")
	screens, err := usecase.GetScreensById(id)
	if err != nil {
		response.Error(c, "Screen not found")
		return
	}

	response.Success(c, "", gin.H{"screens": screens})
}

// UpdateScreen updates a screen
func UpdateScreen(c *gin.Context) {
	id := c.Param("id")
	var screen domain.Screen

	_, err := usecase.GetScreensById(id)
	if err != nil {
		response.Error(c, "Screen not found")
		return
	}

	if err := c.ShouldBindJSON(&screen); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	err = usecase.UpdateScreen(id, screen)
	if err != nil {
		response.Error(c, "Failed to update screen")
		return
	}

	// Send WebSocket message
	SendToClients("Screen updated")

	response.Success(c, "Screen updated successfully", nil)
}

// DeleteScreen deletes a screen
func DeleteScreen(c *gin.Context) {
	id := c.Param("id")

	_, err := usecase.GetScreensById(id)
	if err != nil {
		response.Error(c, "Screen not found")
		return
	}

	// Proceed with deletion if the screen exists
	err = usecase.DeleteScreen(id)
	if err != nil {
		response.Error(c, "Failed to delete screen")
		return
	}

	// Send WebSocket message
	SendToClients("Screen deleted")

	// Return success response
	response.Success(c, "Screen deleted successfully", nil)
}

// ListScreens retrieves all Screens
func ListScreens(c *gin.Context) {
	screen, err := usecase.ListScreens()
	if err != nil {
		response.Error(c, "Failed to fetch Screens")
		return
	}

	response.Success(c, "", gin.H{
		"Screens": screen,
	})
}
