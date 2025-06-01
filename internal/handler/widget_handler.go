package handler

import (
	"github.com/gin-gonic/gin"
	"startfront-backend/internal/domain"
	"startfront-backend/internal/usecase"
	"startfront-backend/pkg/response"
)

// CreateWidget handles the creation of a new widget
func CreateWidget(c *gin.Context) {
	var widget domain.Widget
	if err := c.ShouldBindJSON(&widget); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	err := usecase.CreateWidget(widget)
	if err != nil {
		response.Error(c, "Failed to create widget")
		return
	}

	// Send WebSocket message
	SendToClients("Widget created")

	response.Success(c, "Widget created successfully", nil)
}

// GetWidgetsByScreenID fetches widgets by screen ID
func GetWidgetsByScreenID(c *gin.Context) {
	// screenID := c.Param("screen_id")
	// widgets, err := usecase.GetWidgetsByScreenID(screenID)
	// if err != nil {
	// 	response.Error(c, "Failed to fetch widgets")
	// 	return
	// }

	// response.Success(c, "", gin.H{"widgets": widgets})
}

// UpdateWidget updates widget details
func UpdateWidget(c *gin.Context) {
	// id := c.Param("id")
	var widget domain.Widget
	if err := c.ShouldBindJSON(&widget); err != nil {
		response.Error(c, "Invalid request payload")
		return
	}

	// err := usecase.UpdateWidget(id, widget)
	// if err != nil {
	// 	response.Error(c, "Failed to update widget")
	// 	return
	// }

	// Send WebSocket message
	SendToClients("Widget updated")

	response.Success(c, "Widget updated successfully", nil)
}

// DeleteWidget deletes a widget
func DeleteWidget(c *gin.Context) {
	// id := c.Param("id")
	// err := usecase.DeleteWidget(id)
	// if err != nil {
	// 	response.Error(c, "Failed to delete widget")
	// 	return
	// }

	// Send WebSocket message
	SendToClients("Widget deleted")

	response.Success(c, "Widget deleted successfully", nil)
}
